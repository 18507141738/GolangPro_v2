package controllers

import (
	"Artifice_V2.0.0/models"
)

func placeByOrgainze(organize *models.Organize) []*models.Place {
	o := O
	var maps []*models.Place
	o.QueryTable(new(models.Place)).Filter("Organize__ID", organize.ID).RelatedSel().All(&maps)
	return maps
}
