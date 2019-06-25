package photon_spectator

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/google/gopacket"
)

func TestBasicPhotonLayer(t *testing.T) {
	photonHeader := []byte{
		0x00, 0x01, // PeerID
		0x01,                   // CrcEnabled
		0x01,                   // CommandCount
		0x00, 0x00, 0x00, 0x01, // Timestamp
		0x00, 0x00, 0x00, 0x01, // Challenge
	}

	photonCommand := []byte{
		AcknowledgeType,        // Type
		0x01,                   // ChannelID
		0x01,                   // Flags
		0x04,                   // ReservedByte
		0x00, 0x00, 0x00, 0x0c, // Length
		0x00, 0x00, 0x00, 0x01, // ReliableSequenceNumber
	}

	data := append(photonHeader, photonCommand...)
	packet := gopacket.NewPacket(data, PhotonLayerType, gopacket.Default)

	photonLayer := packet.Layer(PhotonLayerType)

	if photonLayer == nil {
		t.Errorf("Photon layer should be present")
	}

	packetContent, _ := photonLayer.(PhotonLayer)

	if packetContent.PeerID != uint16(1) {
		t.Errorf("PeerID invalid")
	}

	if packetContent.CrcEnabled != uint8(1) {
		t.Errorf("CrcEnabled invalid")
	}

	if packetContent.CommandCount != uint8(1) {
		t.Errorf("CommandCount invalid")
	}

	if packetContent.Timestamp != uint32(1) {
		t.Errorf("Timestamp invalid")
	}

	if packetContent.Challenge != 1 {
		t.Errorf("Challenge invalid")
	}

	if len(packetContent.Commands) != 1 {
		t.Errorf("Commands length invalid")
	}

	command := packetContent.Commands[0]

	if command.Type != AcknowledgeType {
		t.Errorf("Type invalid")
	}

	if command.ChannelID != uint8(1) {
		t.Errorf("ChannelID invalid")
	}

	if command.Flags != uint8(1) {
		t.Errorf("Flags invalid")
	}

	if command.ReservedByte != uint8(4) {
		t.Errorf("ReservedByte invalid")
	}

	if command.Length != PhotonCommandHeaderLength {
		t.Errorf("Length invalid")
	}

	if command.ReliableSequenceNumber != 1 {
		t.Errorf("ReliableSequenceNumber invalid")
	}

	if packetContent.LayerContents() == nil {
		t.Errorf("LayerContents invalid")
	}
}

func TestLoginPhotonLayer(t *testing.T) {
	data, _ := base64.StdEncoding.DecodeString("AAAABLNNdLV46INyBgABAAAAAL0AAAEO8wQBABMAYhABeAAAABDRivQ5Z2OeRKHySYvxxnvBAmI0A3MAJUhJR0hMQU5EX0dSRUVOX01BUktFVFBMQUNFX0NFTlRFUkNJVFkEeQACZr+AAABC2AAACHMABlN5c3RlbQlzAAZTeXN0ZW0NbwEQaQAPQkARbAjW98SV7KgKEmsnEBNsCNb4kxqcYNwUa///FWwI1viTGpxg3BdsCNb4kxqcYNwbbwEcbwEeYgD8awAhBgABAAAAABgAAAEP8wQBAAIAYmT8awEpBgABAAAAAEQAAAEQ8wQBAAkAaQAL+ysBawcsAmIBA3MACU1pc3RpY01hbgRiAQVpBycOAAZ4AAAAAAd4AAAAAPxrABkGAAEAAAADLwAAARHzAwEAACoAQgBpAAv7KgF4AAAAEOPXhPqnUIJLl18ny9qk/KUCcwAJTXlQaWNrbGUzA2IEBngAAAAFAAABAgIHeAAAAAUAAwAAAwhzAAQwMDA3CXkAAmZAwMmOQpe/ugpmQ0PH4AtmRJYAAAxmRJYAAA1mQUAAAA5sCNb5Jscux3MPZkLwAAAQZkLwAAARZj/AAAASbAjW+SbHLsdzE2ZEJ0AAFGZEJ0AAFWZA1hR7FmwI1vkmxy7HcxdmRupgABhmRupgABpsCNV73KAZ8zAcbAAAAACe1CuEHWwAAAAAAHoSAB5sAAAAADt3ipAfbAAAAAAAkPVgIGwI1XvcoBnzMCFsAAAAAlQLvPAibAAAAAAAAAAAJGwI1vkmxy7HcyV5AAVsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACZiBCdzAA92ZXRlcmFuLWZvdW5kZXIocwAPdmV0ZXJhbi1zdGFydGVyKWIAKm8BK3gAAAAQcvpx3aSzQU6SEMkyFwhqZSx5AAppAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAL+ysAAAAAAAAAAC54AAAAEDfDLWHraw5InufShSBq9m8veAAAADAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAxYgA0eQACZgAAAAAAAAAANnMABDAwMDA4eQACZgAAAAAAAAAAOnMAUTk2NzdkMDk2LTc1NDYtNGYzYi1iYmM3LWIwYjQ4OTVhYTcxM0BASVNMQU5EQDBhODY0ZmQxLTM3NTAtNDkyOC1iMDVlLTNmZDhjMGJmNjExZDtsCNb5Jscux3M8abNNdLVEcwAARnMAAEhsAAAAAAEwt9BLcwAATGwI1XvcoBnzME1iHk94AAAAAFB5AABsUXgAAAAAU28BVmwI1vkmoPR0+ldsCNb5JscvPM5YYgFZRGlsAABaeQAFbAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABbbAAAAAAAAAAA/WsAAg==")
	packet := gopacket.NewPacket(data, PhotonLayerType, gopacket.Default)

	photonLayer := packet.Layer(PhotonLayerType)

	if photonLayer == nil {
		t.Errorf("Photon layer should be present")
	}

	packetContent, _ := photonLayer.(PhotonLayer)

	if packetContent.CommandCount != uint8(4) {
		t.Errorf("CommandCount invalid")
	}

	if len(packetContent.Commands) != 4 {
		t.Errorf("Commands length invalid")
	}

	command := packetContent.Commands[3]

	if command.Type != SendReliableType {
		t.Errorf("Type invalid")
	}

	if command.Length != 815 {
		t.Errorf("Length invalid")
	}

	if packetContent.LayerContents() == nil {
		t.Errorf("LayerContents invalid")
	}

	msg, err := command.ReliableMessage()
	if err != nil {
		t.Error(err)
	}

	params := DecodeReliableMessage(msg)
	opCode, ok := params["253"].(int16)
	if params == nil || !ok || opCode != 2 {
		t.Error(err)
	}
	for key, value := range params {
		s, ok := value.(string)
		if ok && strings.HasPrefix(s, "ERROR") {
			t.Errorf("Key %v could not be parsed - %v", key, value)
		}
	}
}

func TestMalformedCommand(t *testing.T) {
	photonHeader := []byte{
		0x00, 0x01, // PeerIdx
		0x01,                   // CrcEnabled
		0x01,                   // CommandCount
		0x00, 0x00, 0x00, 0x01, // Timestamp
		0x00, 0x00, 0x00, 0x01, // Challenge
	}

	photonCommand := []byte{
		AcknowledgeType,        // Type
		0x01,                   // ChannelID
		0x01,                   // Flags
		0x04,                   // ReservedByte
		0x00, 0x0c, 0x0c, 0x0c, // Length
		0x00, 0x00, 0x00, 0x01, // ReliableSequenceNumber
	}

	data := append(photonHeader, photonCommand...)
	packet := gopacket.NewPacket(data, PhotonLayerType, gopacket.Default)

	photonLayer := packet.Layer(PhotonLayerType)

	if photonLayer != nil {
		t.Errorf("Photon layer should be absent")
	}
}
