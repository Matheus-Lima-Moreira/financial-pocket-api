package messages

import goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

var ptBR = []*goi18n.Message{
	{ID: "user.created", Other: "Usuário criado com sucesso"},
	{ID: "user.listed", Other: "Usuários listados com sucesso"},
	{ID: "user.details", Other: "Usuário detalhado com sucesso"},
	{ID: "auth.login_success", Other: "Login realizado com sucesso"},
	{ID: "auth.refresh_success", Other: "Token atualizado com sucesso"},
	{ID: "auth.verify_email_sent", Other: "Email de verificação enviado com sucesso"},
	{ID: "auth.verify_email_success", Other: "Email verificado com sucesso"},
	{ID: "auth.verify_email_send_failed", Other: "Falha ao enviar email de verificação"},
	{ID: "auth.email_not_verified", Other: "Email ainda não verificado"},
	{ID: "auth.resend_verification_email_success", Other: "Email de verificação enviado com sucesso"},
	{ID: "auth.reset_password_email_sent_success", Other: "Email de reset de senha enviado com sucesso"},
	{ID: "auth.reset_password_success", Other: "Senha resetada com sucesso"},
	{ID: "error.missing_body", Other: "Corpo da requisição ausente"},
	{ID: "error.invalid_json", Other: "JSON inválido"},
	{ID: "error.validation", Other: "Erro de validação"},
	{ID: "error.internal", Other: "Erro interno do servidor"},
	{ID: "error.email_already_in_use", Other: "Email já está em uso"},
	{ID: "error.token_already_exists", Other: "Token já existe"},
	{ID: "error.missing_token", Other: "Token ausente"},
	{ID: "error.invalid_token", Other: "Token inválido"},
	{ID: "error.expired_token", Other: "Token expirado"},
	{ID: "action.listed", Other: "Ações listadas com sucesso"},
	{ID: "group_permission.listed", Other: "Permissões de grupo listadas com sucesso"},
	{ID: "auth.reset_password_rate_limited", Other: "Limite de reset de senha atingido"},
}

func GetPTBRMessages() []*goi18n.Message {
	return ptBR
}
