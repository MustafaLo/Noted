package models

type HelpTemplate struct{
	Usage string
	Description string
	Flags []string
	Examples []string
	Notes []string
}