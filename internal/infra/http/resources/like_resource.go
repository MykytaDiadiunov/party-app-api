package resources

type UserLikedDto struct {
	IsLiked bool `json:"isLiked"`
}

func (m UserLikedDto) ResultToDto(result bool) UserLikedDto {
	return UserLikedDto{
		IsLiked: result,
	}
}
