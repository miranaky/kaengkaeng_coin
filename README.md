# kaengkaeng_coin

## wallet

- [x] how to work: signiture, valification, private, public key

  1. [x] hash the msg.
         "I love you" -> hash(x) -> "hashed_message"
  2. [x] generate key pair
         KeyPair (publicK, privateK) (save priv to a file)
  3. [x] sign the hash
         ("hashed_message" + privateK) -> "signature"
  4. [x] verify
         ("hashed_message" + "signature" + publicK)

- [x] persistent wallet : save the wallet and restore the wallet
- [ ] apply to transactions : implement signature, vallification
  - [ ] signed the txID and save txIn.signature
  - [ ] valify the signature when make a new tx
