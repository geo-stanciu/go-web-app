package models

type ResponseModel interface {
	Err() bool
	SErr() string
	Url() string
	SetURL(string)
	HasURL() bool
}

type GenericResponseModel struct {
	BError      bool   `json:"err"`
	SError      string `json:"serr"`
	SUrl        string `json:"-"`
	SSuccessURL string `json:"-"`
	SErrorURL   string `json:"-"`
}

func (r *GenericResponseModel) Err() bool {
	return r.BError
}

func (r *GenericResponseModel) SErr() string {
	return r.SError
}

func (r *GenericResponseModel) getUrl() string {
	if len(r.SUrl) > 0 {
		return r.SUrl
	}

	if r.Err() {
		return r.SErrorURL
	}

	return r.SSuccessURL
}

func (r *GenericResponseModel) Url() string {
	url := r.getUrl()

	if url == "-" {
		return ""
	}

	return url
}

func (r *GenericResponseModel) SetURL(url string) {
	r.SUrl = url
}

func (r *GenericResponseModel) HasURL() bool {
	return len(r.Url()) > 0
}
