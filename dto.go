package main

//UserCreationData ...
type UserCreationData struct {
	Login           string `json:"login,omitempty"`
	Name            string `json:"name,omitempty"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"passwordConfirm,omitempty"`
}

//LoginData ...
type LoginData struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

//CreatePollData ...
type CreatePollData struct {
	Name string `json:"name,omitempty"`
}

//AddOptionData ...
type AddOptionData struct {
	Value string `json:"value,omitempty"`
}

//RemoveOptionData ...
type RemoveOptionData struct {
	Value string `json:"value,omitempty"`
}

//PollVoteData ...
type PollVoteData struct {
	Value string `json:"value,omitempty"`
}

//PollVoteResult ...
type PollVoteResult struct {
	VoteID       string
	VoteCounting map[string]float64
}
