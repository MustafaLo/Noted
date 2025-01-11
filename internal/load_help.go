package internal

import (
	"fmt"
	"strings"
	"github.com/MustafaLo/Noted/models"
)



func GenerateHelpMessage(template models.HelpTemplate)(string){
	return fmt.Sprintf(`
	
Usage:
 %s

Description:
%s
	
Flags:
%s
	
Examples:
%s
	
Notes:
%s
	`, template.Usage, formatParagraph(template.Description), formatList(template.Flags), formatList(template.Examples), formatList(template.Notes))

}

func CreateHelpTemplate(usage string, description string, flags []string, examples []string, notes []string)(models.HelpTemplate){
	return models.HelpTemplate{
		Usage: usage,
		Description: description,
		Flags: flags,
		Examples: examples,
		Notes: notes,
	}
}

func formatParagraph(paragraph string)(string){
	lines := strings.Split(paragraph, ". ")
	var res string
	for _, line := range lines{
		res += " " + line + "\n"
	}
	return res
}

func formatList(items []string)(string){
	if len(items) == 0{
		return "None\n"
	}
	var res string
	for _, item := range items{
		res += " " + item + "\n"
	}
	return res
}