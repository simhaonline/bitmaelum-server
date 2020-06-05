package account

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/rand"
    "crypto/sha256"
    "encoding/json"
    "github.com/jaytaph/mailv2/core"
    "github.com/jaytaph/mailv2/core/config"
    "github.com/sirupsen/logrus"
    "golang.org/x/crypto/pbkdf2"
    "io"
    "io/ioutil"
    "path"
)

type EncryptedAccountInfo struct {
    Address         string  `json:"address"`
    AccountInfo     []byte  `json:"data"`
    Salt            []byte  `json:"salt"`
    Iv              []byte  `json:"iv"`
    Hmac            []byte  `json:"hmac"`
}

const (
    PbkdfIterations = 100002
)

func LoadAccount(addr core.Address, password []byte) (*core.AccountInfo, error) {
    data, err := ioutil.ReadFile(getPath(addr))
    if err != nil {
        return nil, err
    }

    var plainText []byte

    af := &EncryptedAccountInfo{}
    err = json.Unmarshal(data, &af)
    if err != nil || af.AccountInfo == nil {
        logrus.Warning("account file is unprotected with a password.")

        plainText = []byte(data)
    } else {
        // @TODO: Check HMAC

        derivedAESKey := pbkdf2.Key(password, af.Salt, PbkdfIterations, 32, sha256.New)
        aes256, err := aes.NewCipher(derivedAESKey)
        if err != nil {
            return nil, err
        }

        //plainText := make([]byte, len(af.AccountInfo))
        ctr := cipher.NewCTR(aes256, af.Iv)
        ctr.XORKeyStream(af.AccountInfo, plainText)
    }

    ai := &core.AccountInfo{}
    err = json.Unmarshal(plainText, &ai)
    if err != nil {
        return nil, err
    }

    return ai, nil
}

func SaveAccount(addr core.Address, password []byte, acc core.AccountInfo) (error) {
    // Generate JSON structure that we will encrypt
    plainText, err := json.MarshalIndent(&acc, "", "  ")
    if err != nil {
        return err
    }

    // Generate 64 byte salt
    salt := make([]byte, 64)
    _, err = io.ReadFull(rand.Reader, salt)
    if err != nil {
        return err
    }

    // Generate key based on password
    derivedAESKey := pbkdf2.Key(password, salt, PbkdfIterations, 32, sha256.New)
    aes256, err := aes.NewCipher(derivedAESKey)
    if err != nil {
        return err
    }

    // Generate 32 byte IV
    iv := make([]byte, aes.BlockSize)
    _, err = io.ReadFull(rand.Reader, iv)
    if err != nil {
        return err
    }

    // Encrypt the data
    cipherText := make([]byte, len(plainText))
    ctr := cipher.NewCTR(aes256, iv)
    ctr.XORKeyStream(cipherText, plainText)


    // Generate HMAC
    hash := hmac.New(sha256.New, password)
    hash.Write(cipherText)

    af := &EncryptedAccountInfo{
        Address:     addr.String(),
        AccountInfo: cipherText,
        Salt:        salt,
        Iv:          iv,
        Hmac:        hash.Sum(nil),
    }
    data, err := json.MarshalIndent(af, "", "  ")
    if err != nil {
        return err
    }

    // And write to file
    err = ioutil.WriteFile(getPath(addr), data, 0600)
    if err != nil {
        return err
    }

    return nil
}

func getPath(addr core.Address) string {
    return path.Clean(path.Join(config.Client.Accounts.Path, addr.String() + ".account.json"))
}
