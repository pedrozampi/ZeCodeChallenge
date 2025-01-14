package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type autoIncrement struct {
	sync.Mutex
	id int
}

func (a *autoIncrement) ID() (id int) {
	a.Lock()
	defer a.Unlock()

	id = a.id
	a.id++
	return
}

func (a *autoIncrement) SyncID() {
	a.id = Partners[len(Partners)-1].ID
}

var ai autoIncrement

type Partner struct {
	ID           int    `json:"id"`
	TradingName  string `json:"tradingName"`
	OwnerName    string `json:"ownerName"`
	Document     string `json:"document"`
	CoverageArea struct {
		Type        string          `json:"type"`
		Coordinates [][][][]float64 `json:"coordinates"`
	} `json:"coverageArea"`
	Address struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"address"`
}

type Location struct {
	Y           float64 `json:"y"`
	X           float64 `json:"x"`
	MaxDistance float64 `json:"maxDistance"`
}

var Partners = []Partner{
	{
		ID: 0, TradingName: "Adega da Cerveja - Pinheiros",
		OwnerName: "Ze da Silva",
		Document:  "1432132123891/0001",
		CoverageArea: struct {
			Type        string          "json:\"type\""
			Coordinates [][][][]float64 "json:\"coordinates\""
		}{
			Type: "MultiPolygon",
			Coordinates: [][][][]float64{
				{{{30, 20}, {45, 40}, {10, 40}, {30, 20}}},
				{{{15, 5}, {40, 10}, {10, 20}, {5, 10}, {15, 5}}},
			},
		},
		Address: struct {
			Type        string    "json:\"type\""
			Coordinates []float64 "json:\"coordinates\""
		}{
			Type:        "Point",
			Coordinates: []float64{-46.57421, -21.785741},
		},
	},
	{
		ID: 1, TradingName: "Adega do Paiva - Zona Leste",
		OwnerName: "Paiva Silva",
		Document:  "1432132123791/0001",
		CoverageArea: struct {
			Type        string          "json:\"type\""
			Coordinates [][][][]float64 "json:\"coordinates\""
		}{
			Type: "MultiPolygon",
			Coordinates: [][][][]float64{
				{{{-22, 0}, {39, 13}, {-43, -21}, {-117, -13}, {-22, 0}}},
				{{{-85, 53}, {-74, -37}, {-117, -62}, {-51, -83}, {-85, 53}}},
			},
		},
		Address: struct {
			Type        string    "json:\"type\""
			Coordinates []float64 "json:\"coordinates\""
		}{
			Type:        "Point",
			Coordinates: []float64{26.57421, 45.785741},
		},
	},
}

func getPartners(c *gin.Context) {
	c.JSON(http.StatusOK, Partners)
}

func getPartnerByID(c *gin.Context) {
	id := c.GetInt("id")

	for _, p := range Partners {
		if p.ID == id {
			c.JSON(http.StatusFound, p)
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Partner not found."})
}

func registerPartner(c *gin.Context) {
	var newPartner Partner

	if err := c.BindJSON(&newPartner); err != nil {
		return
	}
	ai.SyncID()
	newPartner.ID = ai.ID()

	for _, p := range Partners {
		if p.Document == newPartner.Document {
			c.JSON(http.StatusNotAcceptable, gin.H{"message": "Document already in use."})
			return
		}
	}

	Partners = append(Partners, newPartner)

	c.JSON(http.StatusCreated, newPartner)
}

func nearestPartner(c *gin.Context) {
	var loc Location

	if err := c.BindJSON(&loc); err != nil {
		return
	}
	distanceFromClient := []struct {
		distance float64
		partner  Partner
	}{}
	found := false
	nPartners := len(Partners)
	for k := 0; k < nPartners; k++ {
		for _, p := range Partners[k].CoverageArea.Coordinates {
			for _, l := range p {
				n := len(l)
				intersectCount := 0
				dista := 0.0
				for i := 0; i < n; i++ {

					j := (i + 1) % n
					p1 := l[i]
					p2 := l[j]

					dista = distance([]float64{loc.X, loc.Y}, p1)
					if dista < loc.MaxDistance {
						found = true
					}
					if dista == 0 {
						intersectCount = 3
					}
					if (p1[1] > loc.X) != (p2[1] > loc.Y) {

						xInter := (p2[0]-p1[0])*(loc.Y-p1[1])/(p2[1]-p1[1]) + p1[0]

						if loc.X < xInter {
							intersectCount++
						}
					}

				}
				if intersectCount%2 == 1 {

					distanceFromClient = append(distanceFromClient, struct {
						distance float64
						partner  Partner
					}{
						distance: dista, partner: Partners[k],
					})

				}
			}
		}

	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"message": "No partner found"})
		return
	}
	n := len(distanceFromClient)
	minDist := 0.0
	minPartner := Partner{}

	if n > 0 {
		minDist = distanceFromClient[0].distance
		minPartner = distanceFromClient[0].partner
		for i := 0; i < n; i++ {

			if distanceFromClient[i].distance < minDist {
				minDist = distanceFromClient[i].distance
				minPartner = distanceFromClient[i].partner
			}
		}
		c.JSON(http.StatusFound, gin.H{"Partner": minPartner.TradingName, "Distance": minDist})
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Something went wrong"})
	}

}

func main() {
	router := gin.Default()

	router.GET("/partners", getPartners)
	router.GET("/partners?id=:id", getPartnerByID)

	router.POST("/partners", registerPartner)
	router.POST("/partners/near", nearestPartner)

	router.Run("localhost:8080")
}
