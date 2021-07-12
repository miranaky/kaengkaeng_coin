# kaengkaeng_coin

## wallet

- [ ] how to work: signiture, valification, private, public key

  1. [ ] hash the msg.
         "I love you" -> hash(x) -> "hashed_message"
  2. [ ] generate key pair
         KeyPair (publicK, privateK) (save priv to a file)
  3. [ ] sign the hash
         ("hashed_message" + privateK) -> "signature"
  4. [ ] verify
         ("hashed_message" + "signature" + publicK)

- [ ] persistent wallet : save the wallet and restore the wallet
- [ ] apply to transactions : implement signature, vallification
