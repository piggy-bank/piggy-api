package utils

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"

	"cloud.google.com/go/datastore"
	kms "cloud.google.com/go/kms/apiv1"
	"github.com/manubidegain/piggy-api/cmd/api/configuration"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Entry struct {
	Address    string
	PublicKey  string
	PrivateKey []byte
}

type Readed struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

func SetupDataStoreClient(scope string, projectConfig *configuration.ProjectConfig) (*datastore.Client, error) {
	client, err := datastore.NewClient(context.Background(), projectConfig.ProjectID)
	if err != nil {
		fmt.Printf("Broking inside setup datastore, projectID : %s", projectConfig.ProjectID)
		return client, err
	}
	return client, nil
}

func UploadValue(ctx context.Context, entry Entry, kind string, client *datastore.Client) (string, error) {
	key := datastore.NameKey(kind, entry.Address, nil)
	if key, err := client.Put(ctx, key, &entry); err != nil {
		return key.Name, err
	}
	return key.Name, nil
}

func CreateNewEntry(address string, publicKey string, privateKey string, projectConfig *configuration.ProjectConfig) (Entry, error) {
	entry := Entry{Address: address, PublicKey: publicKey}
	hashedKey, err := EncryptPrivateKey(privateKey, projectConfig)
	if err != nil {
		return entry, err
	}
	entry.PrivateKey = hashedKey
	return entry, nil

}

func EncryptPrivateKey(message string, projectConfig *configuration.ProjectConfig) ([]byte, error) {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// message := "Sample message"

	name := fmt.Sprintf("projects/%s/locations/us-west2/keyRings/addresses-key-ring/cryptoKeys/addresses-key", projectConfig.ProjectID)
	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Convert the message into bytes. Cryptographic plaintexts and
	// ciphertexts are always byte arrays.
	plaintext := []byte(message)

	// Optional but recommended: Compute plaintext's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	plaintextCRC32C := crc32c(plaintext)

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:            name,
		Plaintext:       plaintext,
		PlaintextCrc32C: wrapperspb.Int64(int64(plaintextCRC32C)),
	}

	// Call the API.
	result, err := client.Encrypt(ctx, req)
	if err != nil {
		return nil, err
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if !result.VerifiedPlaintextCrc32C {
		return nil, errors.New("encrypt: request corrupted in-transit")
	}
	if int64(crc32c(result.Ciphertext)) != result.CiphertextCrc32C.Value {
		return nil, errors.New("encrypt: response corrupted in-transit")
	}
	return result.Ciphertext, nil
}

func DecryptPrivateKey(ciphertext []byte, profile string, projectConfig *configuration.ProjectConfig) (string, error) {
	// name := "projects/my-project/locations/us-east1/keyRings/my-key-ring/cryptoKeys/my-key"
	// ciphertext := []byte("...")  // result of a symmetric encryption call
	name := fmt.Sprintf("projects/%s/locations/us-west2/keyRings/addresses-key-ring/cryptoKeys/addresses-key", projectConfig.ProjectID)

	// Create the client.
	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	// Optional, but recommended: Compute ciphertext's CRC32C.
	crc32c := func(data []byte) uint32 {
		t := crc32.MakeTable(crc32.Castagnoli)
		return crc32.Checksum(data, t)
	}
	ciphertextCRC32C := crc32c(ciphertext)

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:             name,
		Ciphertext:       ciphertext,
		CiphertextCrc32C: wrapperspb.Int64(int64(ciphertextCRC32C)),
	}

	// Call the API.
	result, err := client.Decrypt(ctx, req)
	if err != nil {
		return "", err
	}

	// Optional, but recommended: perform integrity verification on result.
	// For more details on ensuring E2E in-transit integrity to and from Cloud KMS visit:
	// https://cloud.google.com/kms/docs/data-integrity-guidelines
	if int64(crc32c(result.Plaintext)) != result.PlaintextCrc32C.Value {
		return "", errors.New("decrypt: response corrupted in-transit")
	}

	return string(result.Plaintext), nil
}

func GetValue(ctx context.Context, client *datastore.Client, stringKey string, kind string, profile string, projectConfig *configuration.ProjectConfig) (Readed, error) {
	key := datastore.NameKey(kind, stringKey, nil)
	entry := Entry{}
	readed := Readed{}
	if err := client.Get(ctx, key, &entry); err != nil {
		fmt.Printf("Client get is broken : %s", projectConfig.ProjectID)
		return readed, err
	}
	readed.Address = entry.Address
	readed.PublicKey = entry.PublicKey
	privateKey, err := DecryptPrivateKey(entry.PrivateKey, profile, projectConfig)
	if err != nil {
		fmt.Printf("Decrypt is broken : %s and profile : %s", projectConfig.ProjectID, profile)
		return readed, err
	}
	readed.PrivateKey = privateKey
	return readed, nil
}
