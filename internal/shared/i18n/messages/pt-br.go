package messages

import goi18n "github.com/nicksnyder/go-i18n/v2/i18n"

var ptBR = []*goi18n.Message{
	{ID: "user.created", Other: "Usuário criado com sucesso"},
	{ID: "user.listed", Other: "Usuários listados com sucesso"},
	{ID: "user.details", Other: "Usuário detalhado com sucesso"},
	{ID: "auth.login_success", Other: "Login realizado com sucesso"},
	{ID: "auth.refresh_success", Other: "Token atualizado com sucesso"},
	{ID: "auth.verify_email_sent", Other: "Email de verificacao enviado com sucesso"},
	{ID: "auth.verify_email_success", Other: "Email verificado com sucesso"},
	{ID: "auth.verify_email_send_failed", Other: "Falha ao enviar email de verificacao"},
	{ID: "auth.email_not_verified", Other: "Email ainda nao verificado"},
	{ID: "error.missing_body", Other: "Corpo da requisição ausente"},
	{ID: "error.invalid_json", Other: "JSON inválido"},
	{ID: "error.validation", Other: "Erro de validação"},
	{ID: "error.internal", Other: "Erro interno do servidor"},
}

func GetPTBRMessages() []*goi18n.Message {
	return ptBR
}
