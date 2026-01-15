package services

import (
	"fmt"
	"strings"

	"smart-choice/repository"
)

type MetaTags struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	OpenGraph   map[string]string `json:"open_graph"`
	Canonical   string            `json:"canonical"`
}

func GetProductMetaTags(productID uint) (*MetaTags, error) {
	product, err := repository.GetProductByID(productID)
	if err != nil {
		return nil, err
	}

	title := fmt.Sprintf("%s - Compre Agora | Smart Choice", product.Name)
	description := fmt.Sprintf("Compre %s por apenas R$%.2f. %s. Frete rápido e seguro.",
		product.Name, product.Price, truncateString(product.Description, 150))

	return &MetaTags{
		Title:       title,
		Description: description,
		OpenGraph: map[string]string{
			"og:title":       title,
			"og:description": description,
			"og:type":        "product",
			"og:image":       fmt.Sprintf("https://smart-choice.com/products/%d/image", product.ID),
			"og:url":         fmt.Sprintf("https://smart-choice.com/products/%d", product.ID),
		},
		Canonical: fmt.Sprintf("https://smart-choice.com/products/%d", product.ID),
	}, nil
}

func GetCategoryMetaTags(category string) *MetaTags {
	title := fmt.Sprintf("%s - Produtos | Smart Choice", strings.Title(category))
	description := fmt.Sprintf("Confira nossa seleção de %s com os melhores preços. Qualidade garantida e entrega rápida.", category)

	return &MetaTags{
		Title:       title,
		Description: description,
		OpenGraph: map[string]string{
			"og:title":       title,
			"og:description": description,
			"og:type":        "website",
			"og:image":       "https://smart-choice.com/images/category-default.jpg",
			"og:url":         fmt.Sprintf("https://smart-choice.com/categories/%s", strings.ToLower(category)),
		},
		Canonical: fmt.Sprintf("https://smart-choice.com/categories/%s", strings.ToLower(category)),
	}
}

func GetHomeMetaTags() *MetaTags {
	title := "Smart Choice - Os Melhores Produtos com os Melhores Preços"
	description := "Descubra uma vasta seleção de produtos de alta qualidade com preços imbatíveis. Compre com segurança e receba em casa."

	return &MetaTags{
		Title:       title,
		Description: description,
		OpenGraph: map[string]string{
			"og:title":       title,
			"og:description": description,
			"og:type":        "website",
			"og:image":       "https://smart-choice.com/images/home-banner.jpg",
			"og:url":         "https://smart-choice.com/",
		},
		Canonical: "https://smart-choice.com/",
	}
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}
