package shodan

type HostSearch struct {
	Matches []struct {
		Info      string        `json:"info,omitempty"`
		Product   string        `json:"product,omitempty"`
		Hostnames []interface{} `json:"hostnames"`
		Hash      int           `json:"hash"`
		IP        int           `json:"ip"`
		Org       string        `json:"org"`
		Isp       string        `json:"isp"`
		Transport string        `json:"transport"`
		Cpe       []string      `json:"cpe,omitempty"`
		Data      string        `json:"data"`
		Asn       string        `json:"asn"`
		Port      int           `json:"port"`
		Version   string        `json:"version,omitempty"`
		IPStr     string        `json:"ip_str,omitempty"`
		Location  struct {
			City         interface{} `json:"city"`
			RegionCode   interface{} `json:"region_code"`
			AreaCode     interface{} `json:"area_code"`
			Longitude    float64     `json:"longitude"`
			CountryCode3 interface{} `json:"country_code3"`
			Latitude     float64     `json:"latitude"`
			PostalCode   interface{} `json:"postal_code"`
			DmaCode      interface{} `json:"dma_code"`
			CountryCode  string      `json:"country_code"`
			CountryName  string      `json:"country_name"`
		} `json:"location"`
	} `json:"matches"`
	Total int `json:"total"`
}
