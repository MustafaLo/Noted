package models

type HelpTemplate struct{
	Heading string
	Usage string
	Description string
	Flags []string
	Examples []string
	Notes []string
}