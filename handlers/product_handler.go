package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/ravik-karn/Dataweave/products"
)

func HandleProduct(logger logrus.Logger, fetcher products.Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := fetcher.Fetch()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("Unable to fetch products: %s", err)
			w.Write([]byte("Unable to fetch products"))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(products)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logger.Errorf("error marshalling response: %s", err)
			w.Write([]byte("Unable to fetch products"))
			return
		}
	}
}
