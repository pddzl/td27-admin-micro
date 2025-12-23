package authority

type UserPlusRoleNameDTO struct {
	AuthorityUserEntity
	RoleName string
}

type UpdateUserDTO struct {
	ID          uint64
	Username    string
	Password    string
	Phone       string
	Email       string
	Active      bool
	RoleModelId uint64
}
