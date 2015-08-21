# Pact
An experimental cryptographic messaging application.

## What Is Pact?
Pact is a CLI application that enables messages to be shared securely between many parties without the 
need for out of band secret sharing. Pact does this by relying on asymmetric cryptography to safeguard a
symmetric key. This allows two parties to communicate securely after exchanging only their public keys.

Pact offloads all of its crypto to the [MSG library](https://github.com/MattSurabian/msg) but this README 
will review the crypto operations as if Pact was doing all the heavy lifting so potential Pact users are given
sufficient context. As of this writing Pact is the only usage of the MSG library.

### Asymmetric and Symmetric Crypto!?
Yes, but they are not used on top of one another. Unlike PGP, Pact does not rely on RSA or DSA public-key crypto. 
Instead, Pact uses [NaCl](http://nacl.cr.yp.to/), a more modern approach to public-key cryptography with 
secure keys that are only 32 bytes long. NaCl's simplicity and security come with a price. Multiple party
decryption of a single cipher text is not possible without shared keys. Instead, new cipher
text must be created for each individual with which a user intends to communicate.

Pact solves this problem by using AES-256-GCM (Galois Counter Mode) to secure the initial message and then
encrypting the secret key used by AES with NaCl. The final cipher text is the concatenation of the AES-256-GCM cipher text
with fixed size repeating blocks of NaCl cipher text containing the key necessary to decrypt the original message. 

### Why Not Just Use PGP?
Frankly, you probably should. This project is an experiment aimed at making NaCl easier to use for the 
"average" person. The only real benefit to Pact is that the keys required for secure communication are 16 
times smaller than the currently recommended 4096-bit RSA keys used for PGP. Pact also aims to be marginally 
easier to use.

### How Does Pact Secure A Message
When Pact encrypts a message it does so using AES-256 in Galois Counter Mode with a randomly generated nonce 
and key. Messages are encrypted for a specific pact, or group of people, which are represented by a collection 
of public keys stored in Pact's configuration file (use `pact list` to see them). Pact loops through these public
keys and encrypts the randomly generated AES-256-GCM key with each pact member's public key. That payload is then
prefixed with the fingerprint of the public key used for encryption, so on decryption the recipient can immediately 
know which chunk of bytes to decrypt first in order to learn the key necessary to decrypt the original message.

### Isn't Combining Cryptographic Method Insecure?
Combining, yes. Concatenating, no. We assume that both AES-256-GCM and NaCl are PRPs(pseudo-random-permutations) 
or at worst PRFs (pseudo-random-functions); which is to say the output they produce is sufficiently indistinguishable 
from actual random output. The concatenation of two pseudo-random blocks is itself pseudo random. All parallelizable 
crypto algorithms rely on this principal. Pact takes advantage of producing a psuedo random block which can be intelligently 
sliced appart by an authorized recipient and securely decrypted.

## Getting Started
There are a couple of ways to get started with Pact. The easiest method is to [download a compiled binary]() for your platform and just start using it. 
If however, you're interested in compiling it yourself the following steps should get you there:

1. Clone this repository
1. Run `make`. This should create a `pact` binary which you can use directly like so: `./pact`; you can also copy the binary to `/usr/bin/` or `/usr/local/bin` to use it directly.
1. Run `pact config` to generate a config file. This will also generate a keypair if one does not already exist, and configure the "self" pact which will contain your own public key.
1. Running `pact list` will show all your pacts and the public keys they contain. A fresh configuration will only have a `self` pact
1. You should now be able to run the create and read methods! Try it: `pact create self "some message" | pact read`

## General Usage


## Available Commands

```
$ ./pact -h
A CLI tool that uses NaCl and AES-256-GCM to facilitate multiparty
communication without the need for out of band secret sharing.
Usage: 
  pact [flags]
  pact [command]

Available Commands: 
  create      Outputs an encrypted ciphertext given a plain text message
  read        Outputs a plain text message given an encrypted cipher text
  config      Configuration generation wizard
  key-gen     Creates new NaCl keys in the location specified by pact's configuration
  key-export  Outputs the user's public key encoded as base64 to STDOUT
  new         Creates a new pact
  rm          Completely removes an existing pact
  list        Lists existing pacts
  add-key     Adds a key to an existing pact or creates a new pact containing the key
  rm-key      Interactively removes a single key from an existing pact
  help        Help about any command

Flags:
  -h, --help=false: help for pact


Use "pact [command] --help" for more information about a command.

```

### create

```
$ ./pact create --help
Uses AES-256-GCM to encrypt a message with a randomly generated key 
from PBKDF2 and encrypts that secret key with the public key of each 
member of a pact. Base64 encoded encrypted ciphertext is sent to STDOUT.

Usage: 
  pact create [pact-name] [plain-text] [flags]

```

### read

```
$ ./pact read --help
Uses NaCl to decrypt a key which can be used to decrypt the message 
which has been secured with AES-256-GCM encryption.

Usage: 
  pact read [cipher-text] [flags]

Flags:
  -h, --help=false: help for read

```

### config

### key-gen

### key-export

### new

### rm

### list

### add-key

### rm-key
