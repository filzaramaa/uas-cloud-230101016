package wa

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"time"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"

	"main/ai"
	"main/models"

	"gorm.io/gorm"
)

// variabel untuk client whatsapp
var clientWa *whatsmeow.Client
var DB *gorm.DB

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())

		fmt.Println(" => dari saya = ", v.Info.IsFromMe)
		fmt.Println(" => server = ", v.Info.MessageSource.Chat.Server)
		fmt.Println(" => apakah group = ", v.Info.IsGroup)
		fmt.Println(" => apakah broadcast = ", v.Info.IsIncomingBroadcast())

		//filter pesan
		if !v.Info.IsFromMe &&
			!v.Info.IsGroup &&
			!v.Info.IsIncomingBroadcast() {

			fmt.Println("PENGIRIM = ", v.Info.Sender.User)
			var pesan string

			switch {
			case v.Message.GetConversation() != "":
				pesan = v.Message.GetConversation()

			case v.Message.GetExtendedTextMessage() != nil:
				pesan = v.Message.GetExtendedTextMessage().GetText()

			case v.Message.GetImageMessage() != nil:
				pesan = v.Message.GetImageMessage().GetCaption()

			case v.Message.GetVideoMessage() != nil:
				pesan = v.Message.GetVideoMessage().GetCaption()
			}

			fmt.Printf("PESAN = [%s]\n", pesan)

			if pesan == "" {
				fmt.Printf("DEBUG MESSAGE = %+v\n", v.Message)
				return
			}
			//membuat array id_wa
			var id_wa []string
			id_wa = append(id_wa, v.Info.ID)

			//status pesan dibaca
			clientWa.MarkRead(context.Background(), id_wa, time.Now(), v.Info.Chat, v.Info.Sender)

			//pengirim akan menerima status
			clientWa.SubscribePresence(context.Background(), v.Info.Sender)

			//status online
			clientWa.SendPresence(context.Background(), types.PresenceAvailable)

			//jeda 3 detik
			time.Sleep(3 * time.Second)

			//status mengetik
			clientWa.SendChatPresence(context.Background(), v.Info.Sender, types.ChatPresenceComposing, types.ChatPresenceMediaText)

			//jeda 3 detik
			time.Sleep(3 * time.Second)

			//status berhenti mengetik
			clientWa.SendChatPresence(context.Background(), v.Info.Sender, types.ChatPresencePaused, types.ChatPresenceMediaText)

			// --- BAGIAN LOGIKA PENGECEKAN PESAN YANG SUDAH DIPERBARUI ---
			pesanCek := strings.ToLower(pesan)

			if pesanCek == "tes" {
				kirimPesan(v.Info.Sender)

			} else {
				// 1. Cek dulu ke database, kalau ada jawaban resmi -> pakai itu
				// 2. Kalau tidak ketemu, teruskan ke AI biar tetap dijawab
				kirimPesanDatabaseAtauAi(v.Info.Sender, pesanCek, pesan)
			}
			// -------------------------------------------------------------
		}
	}
}

func KonekWa(db *gorm.DB) {
	// |------------------------------------------------------------------------------------------------------|
	// | NOTE: You must also import the appropriate DB connector, e.g. github.com/mattn/go-sqlite3 for SQLite |
	// |------------------------------------------------------------------------------------------------------|

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	ctx := context.Background()
	container, err := sqlstore.New(ctx, "sqlite3", "file:wa.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	// If you want multiple sessions, remember their JIDs and use .GetDevice(jid) or .GetAllDevices() instead.
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		panic(err)
	}

	if deviceStore != nil {
		deviceStore.Platform = "macOS"
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	//mengisi variabel client Wa dengan client
	clientWa = client
	DB = db

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	client.Disconnect()
}

func kirimPesan(IDPenerima types.JID) {
	clientWa.SendMessage(context.Background(), IDPenerima, &waE2E.Message{Conversation: proto.String("[UJI COBA] - PESAN OTOMATIS")})
}

// kirimPesanDatabaseAtauAi: cek database dulu untuk jawaban resmi/baku.
// Kalau kode tidak ditemukan, pertanyaan diteruskan ke AI supaya tetap dijawab.
func kirimPesanDatabaseAtauAi(IDPenerima types.JID, kode string, pesanAsli string) {
	var pesanDb models.Pesan
	result := DB.Where("kode = ?", kode).First(&pesanDb)

	if result.Error == nil {
		// Ketemu jawaban resmi di database
		kirimPesanText(IDPenerima, pesanDb.Balasan)
		return
	}

	// Tidak ketemu di database, teruskan ke AI
	fmt.Println("DEBUG: kode tidak ditemukan di DB, diteruskan ke AI ->", kode)
	userID := IDPenerima.User
	jawabanAi := ai.TanyaAi(userID, pesanAsli)
	kirimPesanText(IDPenerima, jawabanAi)
}

func kirimPesanText(IDPenerima types.JID, isiPesan string) {
	clientWa.SendMessage(
		context.Background(),
		IDPenerima,
		&waE2E.Message{
			Conversation: proto.String(isiPesan),
		},
	)
}
