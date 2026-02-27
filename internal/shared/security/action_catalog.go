package security

type ActionDefinition struct {
	ID          string
	Resource    string
	Action      string
	Label       string
	Description string
}

const (
	DefaultActionCreate  = "create"
	DefaultActionRead    = "read"
	DefaultActionUpdate  = "update"
	DefaultActionDelete  = "delete"
	DefaultActionList    = "list"
	DefaultActionDetails = "details"
)

const (
	ResourceActions         = "actions"
	ResourceUsers           = "users"
	ResourceOrganizations   = "organizations"
	ResourceGroupPermission = "group-permissions"
)

const (
	ActionActionsList = "actions:list"

	ActionUsersList    = "users:list"
	ActionUsersDetails = "users:details"

	ActionOrganizationsList    = "organizations:list"
	ActionOrganizationsDetails = "organizations:details"
	ActionOrganizationsCreate  = "organizations:create"
	ActionOrganizationsUpdate  = "organizations:update"
	ActionOrganizationsDelete  = "organizations:delete"

	ActionGroupPermissionsList    = "group-permissions:list"
	ActionGroupPermissionsDetails = "group-permissions:details"
	ActionGroupPermissionsCreate  = "group-permissions:create"
	ActionGroupPermissionsUpdate  = "group-permissions:update"
	ActionGroupPermissionsDelete  = "group-permissions:delete"
)

var ActionCatalog = []ActionDefinition{
	{
		ID:          ActionActionsList,
		Resource:    ResourceActions,
		Action:      DefaultActionList,
		Label:       "Listar Acoes",
		Description: "Permite listar as acoes do sistema.",
	},
	{
		ID:          ActionUsersList,
		Resource:    ResourceUsers,
		Action:      DefaultActionList,
		Label:       "Listar Usuarios",
		Description: "Permite listar usuarios.",
	},
	{
		ID:          ActionUsersDetails,
		Resource:    ResourceUsers,
		Action:      DefaultActionDetails,
		Label:       "Detalhar Usuario",
		Description: "Permite visualizar detalhes de um usuario.",
	},
	{
		ID:          ActionOrganizationsList,
		Resource:    ResourceOrganizations,
		Action:      DefaultActionList,
		Label:       "Listar Organizacoes",
		Description: "Permite listar organizacoes.",
	},
	{
		ID:          ActionOrganizationsDetails,
		Resource:    ResourceOrganizations,
		Action:      DefaultActionDetails,
		Label:       "Detalhar Organizacao",
		Description: "Permite visualizar detalhes de uma organizacao.",
	},
	{
		ID:          ActionOrganizationsCreate,
		Resource:    ResourceOrganizations,
		Action:      DefaultActionCreate,
		Label:       "Criar Organizacao",
		Description: "Permite criar organizacoes.",
	},
	{
		ID:          ActionOrganizationsUpdate,
		Resource:    ResourceOrganizations,
		Action:      DefaultActionUpdate,
		Label:       "Atualizar Organizacao",
		Description: "Permite atualizar organizacoes.",
	},
	{
		ID:          ActionOrganizationsDelete,
		Resource:    ResourceOrganizations,
		Action:      DefaultActionDelete,
		Label:       "Excluir Organizacao",
		Description: "Permite excluir organizacoes.",
	},
	{
		ID:          ActionGroupPermissionsList,
		Resource:    ResourceGroupPermission,
		Action:      DefaultActionList,
		Label:       "Listar Grupos de Permissao",
		Description: "Permite listar grupos de permissao.",
	},
	{
		ID:          ActionGroupPermissionsDetails,
		Resource:    ResourceGroupPermission,
		Action:      DefaultActionDetails,
		Label:       "Detalhar Grupo de Permissao",
		Description: "Permite visualizar detalhes de um grupo de permissao.",
	},
	{
		ID:          ActionGroupPermissionsCreate,
		Resource:    ResourceGroupPermission,
		Action:      DefaultActionCreate,
		Label:       "Criar Grupo de Permissao",
		Description: "Permite criar grupos de permissao.",
	},
	{
		ID:          ActionGroupPermissionsUpdate,
		Resource:    ResourceGroupPermission,
		Action:      DefaultActionUpdate,
		Label:       "Atualizar Grupo de Permissao",
		Description: "Permite atualizar grupos de permissao.",
	},
	{
		ID:          ActionGroupPermissionsDelete,
		Resource:    ResourceGroupPermission,
		Action:      DefaultActionDelete,
		Label:       "Excluir Grupo de Permissao",
		Description: "Permite excluir grupos de permissao.",
	},
}
