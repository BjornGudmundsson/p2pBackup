package utilities

import (
	"flag"
	"fmt"
	"github.com/BjornGudmundsson/p2pBackup/kyber/group/curve25519"
)

const testRepo = "https://github.com/BjornGudmundsson/TestRepo"

type ErrorNotFound struct {
	flag string
}

func NewErrorNotFound(s string) *ErrorNotFound {
	e := new(ErrorNotFound)
	e.flag = s
	return e
}

func (e *ErrorNotFound) Error() string {
	return fmt.Sprintf("Error: flag %s not found", e.flag)
}

type FlagsContainer struct{
	ints map[string]int
	strings map[string]string
	booleans map[string]bool
}

type Flags interface{
	GetString(s string) string
	GetInt(s string) int
	GetBool(s string) bool
}


func (f *FlagsContainer) GetString(flag string) string {
	if v, ok := f.strings[flag]; ok {
		return v
	}
	return ""
}

func (f * FlagsContainer) GetInt(flag string) int {
	if v, ok := f.ints[flag]; ok {
		return v
	}
	return 0
}

func (f *FlagsContainer) GetBool(flag string) bool {
	if v, ok := f.booleans[flag]; ok {
		return v
	}
	return false
}

//Flags this is to simplify having a certain set of flags and sending them between functions
//without having to call flag.Parse repeatedly
func NewFlags() Flags {
	port := flag.String("p", "8080", "Which port to run the server on")
	baseDir := flag.String("base", ".", "Base is the base directory in which all files will be backed up from. If not provided it will default to the running directory")
	peersList := flag.String("peers", "peers.txt", "Peers is the file in which the data about other peers is stored")
	udpPort := flag.String("udp", "5000", "UDP is the port that will be used for the udp socket")
	rules := flag.String("backuprules", "", "backuprules is the toml file in which the specifications for the backup are kept")
	gui := flag.Bool("gui", false, "Gui says whether or not a gui should be displayed or not. Defaults to false")
	storageFile := flag.String("storage", "backup.txt", "Storage is the location in which you prefer to store your peers backups")
	filePort := flag.Int("fileport", 3000, "The port in which a tcp connection can be made to send the backup")
	backupLogs := flag.String("logfile", "backuplog.txt", "This is where the users wishes to store all log of backups they have performed")
	updateTimer := flag.String("backuprate", "1s", "Backuprate tells how often the system should scan for whether it should update")
	initialize := flag.Bool("init", false, "Init is to tell wheter the user wants to get a new private/public key pair")
	key := flag.String("key", "", "What the secret key that the user will be using to encrypt their backups")
	authString := flag.String("authkey", "", "The key used to prove you are part of the p2p backup group")
	suite := flag.String("suite", curve25519.NewBlakeSHA256Curve25519(true).String(), "What ciphersuite the user decides to use")
	suiteFile := flag.String("Suites", "", "Where all the known suites are stored in a TOML file")
	setString := flag.String("set", "", "The file containing the public keys of the anonymity set")
	retrieve := flag.Bool("retrieve", false, "Retrieving or storing back ups")
	retrievalPassword := flag.String("pw", "deadbeef", "Password used for encrypting/decrypting backups")
	protocol := flag.String("protocol", "tcp", "What protocol should be used")
	findProtocol := flag.String("find", "udp", "how to find other peers")
	update := flag.Bool("update", true, "whether it should be periodically backing up data")
	repo := flag.String("repo", testRepo, "If the git protocol is used, where the peer list is kept")
	ip := flag.String("ip", "127.0.0.1", "The IP address of this peer")
	flag.Parse()
	flags := &FlagsContainer{
		ints: make(map[string]int),
		strings: make(map[string]string),
		booleans: make(map[string]bool),
	}
	flags.strings["p"] = *port
	flags.strings["base"] = *baseDir
	flags.strings["peers"] = *peersList
	flags.strings["udp"] = *udpPort
	flags.strings["backuprules"] = *rules
	flags.strings["storage"] = *storageFile
	flags.strings["logfile"] = *backupLogs
	flags.strings["backuprate"] = *updateTimer
	flags.strings["key"] = *key
	flags.strings["authkey"] = *authString
	flags.strings["suite"] = *suite
	flags.strings["suites"] = *suiteFile
	flags.strings["set"] = *setString
	flags.strings["pw"] = *retrievalPassword
	flags.booleans["gui"] = *gui
	flags.booleans["retrieve"] = *retrieve
	flags.booleans["init"] = *initialize
	flags.ints["fileport"] = *filePort
	flags.strings["protocol"] = *protocol
	flags.strings["find"] = *findProtocol
	flags.booleans["update"] = *update
	flags.strings["repo"] = *repo
	flags.strings["ip"] = *ip
	return flags
}

