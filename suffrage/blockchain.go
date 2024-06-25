package suffrage

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

//RegisterVote registers a vote in our blockchain
func (b *Blockchain) RegisterVote(vote Vote) bool {
	vote.VoterName = strings.ToLower(vote.VoterName)
	vote.CandidateID = strings.ToLower(vote.CandidateID)
	vote.CandidateName = strings.ToLower(vote.CandidateName)
	b.PendingVotes = append(b.PendingVotes, vote)
	return true
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//RegisterNode registers a node in our blockchain
func (b *Blockchain) RegisterNode(node string) bool {
	if !contains(b.NetworkNodes, node) {
		b.NetworkNodes = append(b.NetworkNodes, node)
	}
	return true
}

//CreateNewBlock ...
func (b *Blockchain) CreateNewBlock(nonce int, previousBlockHash string, hash string) Block {
	newBlock := Block{
		Index:     len(b.Chain) + 1,
		Votes:     b.PendingVotes,
		Timestamp: time.Now().UnixNano(),
		Nonce:     nonce,
		Hash:      hash, PreviousBlockHash: previousBlockHash}

	b.PendingVotes = Votes{}
	b.Chain = append(b.Chain, newBlock)
	return newBlock
}

//GetLastBlock ...
func (b *Blockchain) GetLastBlock() Block {
	return b.Chain[len(b.Chain)-1]
}

//HashBlock ...
func (b *Blockchain) HashBlock(previousBlockHash string, currentBlockData string, nonce int) string {
	h := sha256.New()
	strToHash := previousBlockHash + currentBlockData + strconv.Itoa(nonce)
	h.Write([]byte(strToHash))
	hashed := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hashed
}

//ProofOfWork ...
func (b *Blockchain) ProofOfWork(previousBlockHash string, currentBlockData string) int {
	nonce := -1
	inputFmt := ""
	for inputFmt != "0000" {
		nonce = nonce + 1
		hash := b.HashBlock(previousBlockHash, currentBlockData, nonce)
		inputFmt = hash[0:4]
	}
	return nonce
}

//CheckNewBlockHash ...
func (b *Blockchain) CheckNewBlockHash(newBlock Block) bool {
	lastBlock := b.GetLastBlock()
	correctHash := lastBlock.Hash == newBlock.PreviousBlockHash
	correctIndex := (lastBlock.Index + 1) == newBlock.Index

	return (correctHash && correctIndex)
}

//ChainIsValid Used by consensus algorithm
func (b *Blockchain) ChainIsValid() bool {
	i := 1
	for i < len(b.Chain) {
		currentBlock := b.Chain[i]
		prevBlock := b.Chain[i-1]
		currentBlockData := BlockData{Index: strconv.Itoa(prevBlock.Index - 1), Votes: currentBlock.Votes}
		currentBlockDataAsByteArray, _ := json.Marshal(currentBlockData)
		currentBlockDataAsStr := base64.URLEncoding.EncodeToString(currentBlockDataAsByteArray)
		blockHash := b.HashBlock(prevBlock.Hash, currentBlockDataAsStr, currentBlock.Nonce)

		if blockHash[0:4] != "0000" {
			return false
		}

		if currentBlock.PreviousBlockHash != prevBlock.Hash {
			return false
		}

		i = i + 1
	}

	genesisBlock := b.Chain[0]
	correctNonce := genesisBlock.Nonce == 100
	correctPreviousBlockHash := genesisBlock.PreviousBlockHash == "0"
	correctHash := genesisBlock.Hash == "0"
	correctVotes := len(genesisBlock.Votes) == 0

	return (correctNonce && correctPreviousBlockHash && correctHash && correctVotes)
}

//GetVotesForCandidate ...
func (b *Blockchain) GetVotesForCandidate(candidateID string) Votes {
	candidateVotes := Votes{}
	i := 0
	chainLength := len(b.Chain)
	for i < chainLength {
		block := b.Chain[i]
		votesInBlock := block.Votes
		j := 0
		votesLength := len(votesInBlock)
		for j < votesLength {
			vote := votesInBlock[j]
			if vote.CandidateID == candidateID {
				candidateVotes = append(candidateVotes, vote)
			}
			j = j + 1
		}
		i = i + 1
	}
	return candidateVotes
}

//GetVotesForVoter ...
func (b *Blockchain) GetVotesForVoter(voterName string) Votes {
	candidateVotes := Votes{}
	i := 0
	chainLength := len(b.Chain)
	for i < chainLength {
		block := b.Chain[i]
		votesInBlock := block.Votes
		j := 0
		votesLength := len(votesInBlock)
		for j < votesLength {
			vote := votesInBlock[j]
			if vote.VoterName == voterName {
				candidateVotes = append(candidateVotes, vote)
			}
			j = j + 1
		}
		i = i + 1
	}
	return candidateVotes
}
