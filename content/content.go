package content

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Financial-Times/message-queue-gonsumer/consumer"
	"github.com/satori/go.uuid"
)

// Content is the interface for different type of contents from different CMSs.
type Content interface {
	IsValid() bool
	IsMarkedDeleted() bool
	GetType() string
	GetUUID() string
}

const systemIDKey = "Origin-System-Id"

const dateLayout = "2006-01-02T15:04:05.000Z"
const logPattern = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile | log.LUTC

var info *log.Logger
var warn *log.Logger
var errorL *log.Logger

func init() {
	//to be used for INFO-level logging: info.Println("foo is now bar")
	info = log.New(os.Stdout, "INFO  - ", logPattern)
	//to be used for WARN-level logging: warn.Println("foo is now bar")
	warn = log.New(os.Stdout, "WARN  - ", logPattern)
	//to be used for ERROR-leve logging: errorL.Println("foo is now bar")
	errorL = log.New(os.Stdout, "ERROR - ", logPattern)
}

// UnmarshalContent unmarshals the message body into the appropriate content type based on the systemID header.
func UnmarshalContent(msg consumer.Message) (Content, error) {
	headers := msg.Headers
	systemID := headers[systemIDKey]
	switch systemID {
	case "http://cmdb.ft.com/systems/methode-web-pub":
		var eomFile EomFile
		err := json.Unmarshal([]byte(msg.Body), &eomFile)
		return eomFile, err
	case "http://cmdb.ft.com/systems/wordpress":
		var wordPressMsg WordPressMessage
		err := json.Unmarshal([]byte(msg.Body), &wordPressMsg)
		return wordPressMsg, err
	default:
		return nil, fmt.Errorf("Unsupported content with system ID: [%s].", systemID)
	}
}

func isUUIDValid(contentUUID string) bool {
	parsedUUID, err := uuid.FromString(contentUUID)
	if err != nil {
		warn.Printf("Cannot parse UUID [%v], error: [%v]", contentUUID, err.Error())
		return false
	}
	return contentUUID == parsedUUID.String()
}
