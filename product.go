package main

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func createProductDB(db *gorm.DB, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("ราคาสินค้าต้องไม่ติดลบ (price cannot be negative)")
	}
	if quantity < 0 {
		return nil, errors.New("จำนวนสต็อกต้องไม่ติดลบ (quantity cannot be negative)")
	}

	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("ไม่พบประเภทสินค้า ID %d: %w", categoryID, err)
	}

	product := &Product{
		Name:       name,
		Price:      price,
		Quantity:   quantity,
		CategoryID: categoryID,
	}

	if err := db.Create(product).Error; err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	db.Preload("Category").First(product, product.ID)
	return product, nil
}

func getAllProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch all products: %w", err)
	}
	return products, nil
}

func getProductByIDDB(db *gorm.DB, id uint) (*Product, error) {
	var product Product
	if err := db.Preload("Category").First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}
	return &product, nil
}

func getProductByNameDB(db *gorm.DB, keyword string) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("name LIKE ?", "%"+keyword+"%").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to find products by name keyword '%s': %w", keyword, err)
	}
	return products, nil
}

func getProductsByCategoryDB(db *gorm.DB, categoryID uint) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products by category %d: %w", categoryID, err)
	}
	return products, nil
}

func getProductsByPriceRangeDB(db *gorm.DB, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("price >= ? AND price <= ?", minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products in price range [%.2f, %.2f]: %w", minPrice, maxPrice, err)
	}
	return products, nil
}

func getProductsLowStockDB(db *gorm.DB, threshold int) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("quantity < ?", threshold).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch low stock products (< %d): %w", threshold, err)
	}
	return products, nil
}

func getProductsByCategoryAndPriceDB(db *gorm.DB, categoryID uint, minPrice float64, maxPrice float64) ([]Product, error) {
	var products []Product
	if err := db.Preload("Category").Where("category_id = ? AND price >= ? AND price <= ?", categoryID, minPrice, maxPrice).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products by category %d and price range [%.2f, %.2f]: %w", categoryID, minPrice, maxPrice, err)
	}
	return products, nil
}

func getProductsSortedByPriceDB(db *gorm.DB, sortOrder string) ([]Product, error) {
	upperOrder := strings.ToUpper(strings.TrimSpace(sortOrder))
	if upperOrder != "ASC" && upperOrder != "DESC" {
		upperOrder = "ASC"
	}

	var products []Product
	if err := db.Preload("Category").Order("price " + upperOrder).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sorted products: %w", err)
	}
	return products, nil
}

func getTopExpensiveProductsDB(db *gorm.DB, n int) ([]Product, error) {
	if n <= 0 {
		n = 3
	}
	var products []Product
	if err := db.Preload("Category").Order("price DESC").Limit(n).Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch top %d expensive products: %w", n, err)
	}
	return products, nil
}

func updateProductByIDDB(db *gorm.DB, id uint, name string, price float64, quantity int, categoryID uint) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	if price < 0 {
		return nil, errors.New("ราคาสินค้าต้องไม่ติดลบ")
	}
	if quantity < 0 {
		return nil, errors.New("สต็อกสินค้าต้องไม่ติดลบ")
	}

	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		return nil, fmt.Errorf("category ID %d not found: %w", categoryID, err)
	}

	product.Name = name
	product.Price = price
	product.Quantity = quantity
	product.CategoryID = categoryID

	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product ID %d: %w", id, err)
	}

	db.Preload("Category").First(&product, product.ID)
	return &product, nil
}

func updateProductPriceDB(db *gorm.DB, id uint, newPrice float64) (*Product, error) {
	if newPrice < 0 {
		return nil, errors.New("ราคาต้องไม่ติดลบ (price cannot be negative)")
	}

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}

	product.Price = newPrice
	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update product price: %w", err)
	}

	return &product, nil
}

func updateProductStockDB(db *gorm.DB, id uint, quantityChange int) (*Product, error) {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return nil, fmt.Errorf("product ID %d not found: %w", id, err)
	}

	newStock := product.Quantity + quantityChange
	if newStock < 0 {
		return nil, fmt.Errorf("สต็อกสินค้าต้องไม่ติดลบ (สต็อกเดิม: %d, เปลี่ยนแปลง: %d, สต็อกใหม่: %d)", product.Quantity, quantityChange, newStock)
	}

	product.Quantity = newStock
	if err := db.Save(&product).Error; err != nil {
		return nil, fmt.Errorf("failed to update stock for product ID %d: %w", id, err)
	}

	return &product, nil
}

func softDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	if err := db.Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to soft delete product ID %d: %w", id, err)
	}
	return nil
}

func hardDeleteProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		return fmt.Errorf("product ID %d not found: %w", id, err)
	}

	var orderCount int64
	if err := db.Model(&Order{}).Where("product_id = ?", id).Count(&orderCount).Error; err != nil {
		return fmt.Errorf("failed to check orders for product ID %d: %w", id, err)
	}

	if orderCount > 0 {
		return fmt.Errorf("ห้ามลบสินค้าถาวร: สินค้า ID %d มีประวัติคำสั่งซื้ออยู่ในระบบ %d รายการ", id, orderCount)
	}

	if err := db.Unscoped().Delete(&product).Error; err != nil {
		return fmt.Errorf("failed to hard delete product ID %d: %w", id, err)
	}
	return nil
}

func getDeletedProductsDB(db *gorm.DB) ([]Product, error) {
	var products []Product
	if err := db.Unscoped().Where("deleted_at IS NOT NULL").Preload("Category").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to get soft deleted products: %w", err)
	}
	return products, nil
}

func restoreProductDB(db *gorm.DB, id uint) error {
	var product Product
	if err := db.Unscoped().Where("id = ? AND deleted_at IS NOT NULL", id).First(&product).Error; err != nil {
		return fmt.Errorf("soft deleted product ID %d not found: %w", id, err)
	}

	if err := db.Unscoped().Model(&product).Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("failed to restore product ID %d: %w", id, err)
	}
	return nil
}
