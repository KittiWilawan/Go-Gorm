package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

func printDivider(title string) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Printf("   %s\n", title)
	fmt.Println(strings.Repeat("=", 70))
}

func printSubHeader(subtitle string) {
	fmt.Printf("\n--- [ %s ] ---\n", subtitle)
}

func seedDatabase(db *gorm.DB) {
	var catCount int64
	db.Model(&Category{}).Count(&catCount)
	if catCount > 0 {
		fmt.Printf("ℹ️ ฐานข้อมูลมีข้อมูลอยู่แล้ว (%d categories) ข้ามขั้นตอนการ Seed\n", catCount)
		return
	}

	fmt.Println("🌱 กำลังเริ่มต้นใส่ข้อมูลทดสอบขั้นต่ำ (Seeding data)...")

	categoriesData := []struct {
		Name string
		Desc string
	}{
		{"Laptops & Computers", "คอมพิวเตอร์พกพา แล็ปท็อป และคอมพิวเตอร์ตั้งโต๊ะ"},
		{"Smartphones & Tablets", "สมาร์ทโฟน แท็บเล็ต และอุปกรณ์สื่อสารเคลื่อนที่"},
		{"Accessories & Peripherals", "เมาส์ คีย์บอร์ด ดองเกิล และอุปกรณ์ต่อพ่วงคอมพิวเตอร์"},
		{"Audio & Gaming Gear", "หูฟัง ลำโพง และอุปกรณ์สำหรับเกมเมอร์"},
	}

	catMap := make(map[int]uint)
	for i, c := range categoriesData {
		cat, err := createCategoryDB(db, c.Name, c.Desc)
		if err != nil {
			log.Fatalf("Seed category failed: %v", err)
		}
		catMap[i+1] = cat.ID
	}
	fmt.Printf("✅ เพิ่มหมวดหมู่สำเร็จ: %d รายการ\n", len(catMap))

	customersData := []struct {
		Name  string
		Email string
		Phone string
	}{
		{"สมชาย ใจดี", "somchai.j@example.com", "081-111-1111"},
		{"สุดา รักดี", "suda.r@example.com", "082-222-2222"},
		{"กิตติ วิลาวรรณ", "kitti.w@example.com", "083-333-3333"},
		{"อนันต์ ปัญญาพร", "anan.p@example.com", "084-444-4444"},
		{"มาลี บุญมา", "malee.b@example.com", "085-555-5555"},
		{"วิชัย ชัยสุข", "wichai.c@example.com", "086-666-6666"},
		{"ณัฐวุฒิ วงศ์เทพ", "natthawut.w@example.com", "087-777-7777"},
		{"พรทิพย์ สุขสำราญ", "pornthip.s@example.com", "088-888-8888"},
	}

	custMap := make(map[int]uint)
	for i, c := range customersData {
		cust, err := createCustomerDB(db, c.Name, c.Email, c.Phone)
		if err != nil {
			log.Fatalf("Seed customer failed: %v", err)
		}
		custMap[i+1] = cust.ID
	}
	fmt.Printf("✅ เพิ่มลูกค้าสำเร็จ: %d รายการ\n", len(custMap))

	productsData := []struct {
		Name     string
		Price    float64
		Quantity int
		CatKey   int
	}{
		{"MacBook Air M3 13-inch", 39900.00, 15, 1},
		{"ASUS Vivobook 15 OLED", 21900.00, 20, 1},
		{"Lenovo ThinkPad E14 Gen 5", 28900.00, 10, 1},
		{"Dell XPS 13 Plus", 52900.00, 8, 1},
		{"iPhone 16 128GB", 29900.00, 25, 2},
		{"iPhone 16 Pro Max 256GB", 48900.00, 12, 2},
		{"Samsung Galaxy S24 Ultra", 43900.00, 18, 2},
		{"Samsung Galaxy Tab S9", 25900.00, 14, 2},
		{"iPad Air 11-inch M2", 23900.00, 16, 2},
		{"Logitech Mouse MX Master 3S", 3890.00, 30, 3},
		{"Mechanical Keyboard Keychron K2 Pro", 3490.00, 4, 3},
		{"USB-C Hub 8-in-1 Baseus Metal Gleam", 1290.00, 45, 3},
		{"Anker 65W GaN Fast Charger", 990.00, 3, 3},
		{"Sony WH-1000XM5 Wireless Headphones", 13900.00, 9, 4},
		{"HyperX Cloud III Gaming Headset", 3290.00, 15, 4},
	}

	prodMap := make(map[int]uint)
	for i, p := range productsData {
		prod, err := createProductDB(db, p.Name, p.Price, p.Quantity, catMap[p.CatKey])
		if err != nil {
			log.Fatalf("Seed product failed: %v", err)
		}
		prodMap[i+1] = prod.ID
	}
	fmt.Printf("✅ เพิ่มสินค้าสำเร็จ: %d รายการ\n", len(prodMap))

	ordersData := []struct {
		CustKey int
		ProdKey int
		Qty     int
		Status  string
	}{
		{1, 1, 1, "Completed"},
		{2, 5, 1, "Completed"},
		{3, 10, 2, "Pending"},
		{4, 11, 1, "Completed"},
		{5, 12, 3, "Pending"},
		{6, 7, 1, "Completed"},
		{7, 14, 1, "Pending"},
		{8, 2, 1, "Completed"},
		{1, 13, 1, "Cancelled"},
		{2, 10, 1, "Completed"},
		{3, 6, 1, "Pending"},
		{4, 15, 2, "Completed"},
	}

	for _, o := range ordersData {
		ord, err := createOrderDB(db, custMap[o.CustKey], prodMap[o.ProdKey], o.Qty)
		if err != nil {
			log.Fatalf("Seed order failed: %v", err)
		}
		if o.Status != "Pending" {
			_, err = updateOrderStatusDB(db, ord.ID, o.Status)
			if err != nil {
				log.Fatalf("Update seed order status failed: %v", err)
			}
		}
	}
	fmt.Printf("✅ เพิ่มคำสั่งซื้อสำเร็จ: %d รายการ (ผ่าน Transaction และตัดสต็อกอัตโนมัติ)\n", len(ordersData))
}

func main() {
	printDivider("ระบบจัดการร้านค้า (Store Management System) - Go + GORM + MySQL")

	printSubHeader("ส่วนที่ 1: Database & Models")
	db, err := InitDB()
	if err != nil {
		log.Fatalf("เชื่อมต่อฐานข้อมูลล้มเหลว: %v", err)
	}
	fmt.Println("✅ เชื่อมต่อฐานข้อมูล store_db และสร้างตารางด้วย AutoMigrate() สำเร็จ")

	seedDatabase(db)

	printSubHeader("ส่วนที่ 2: Create (ทดสอบฟังก์ชันการสร้างข้อมูล)")

	newCat, err := createCategoryDB(db, "Smart Home & IoT", "อุปกรณ์สมาร์ทโฮมและกล้องวงจรปิด")
	if err != nil {
		fmt.Printf("❌ สร้าง Category ผิดพลาด: %v\n", err)
	} else {
		fmt.Printf("✅ createCategoryDB(): เพิ่มประเภท ID: %d, Name: %s\n", newCat.ID, newCat.Name)
	}

	newProd, err := createProductDB(db, "Xiaomi Smart Camera C400", 1490.00, 10, newCat.ID)
	if err != nil {
		fmt.Printf("❌ สร้าง Product ผิดพลาด: %v\n", err)
	} else {
		fmt.Printf("✅ createProductDB(): เพิ่มสินค้า ID: %d, Name: %s, Price: %.2f, Category: %s\n",
			newProd.ID, newProd.Name, newProd.Price, newProd.Category.Name)
	}

	testEmail := fmt.Sprintf("test.user.%d@example.com", time.Now().UnixNano())
	newCust, err := createCustomerDB(db, "ทดสอบ ผู้ใช้", testEmail, "089-999-9999")
	if err != nil {
		fmt.Printf("❌ สร้าง Customer ผิดพลาด: %v\n", err)
	} else {
		fmt.Printf("✅ createCustomerDB(): เพิ่มลูกค้า ID: %d, Name: %s, Email: %s\n",
			newCust.ID, newCust.Name, newCust.Email)
	}

	fmt.Println("🛒 ทดสอบ createOrderDB() ปกติ (สั่งซื้อ 2 ชิ้น):")
	newOrder, err := createOrderDB(db, newCust.ID, newProd.ID, 2)
	if err != nil {
		fmt.Printf("❌ สั่งซื้อผิดพลาด: %v\n", err)
	} else {
		fmt.Printf("✅ createOrderDB(): สำเร็จ! Order ID: %d, ลูกค้า: %s, สินค้า: %s, จำนวน: %d, ยอดรวม: %.2f, สถานะ: %s\n",
			newOrder.ID, newOrder.Customer.Name, newOrder.Product.Name, newOrder.Quantity, newOrder.TotalPrice, newOrder.Status)
	}

	fmt.Println("🛒 ทดสอบ createOrderDB() เมื่อสต็อกไม่เพียงพอ (สั่งซื้อ 999 ชิ้น):")
	_, errRollback := createOrderDB(db, newCust.ID, newProd.ID, 999)
	if errRollback != nil {
		fmt.Printf("✅ Transaction ปลอดภัย Rollback ตามเงื่อนไข: %v\n", errRollback)
	}

	printSubHeader("ส่วนที่ 3: Read / Conditional Queries")

	allProducts, _ := getAllProductsDB(db)
	fmt.Printf("1. getAllProductsDB(): มีสินค้าทั้งหมด %d รายการ\n", len(allProducts))

	prod1, _ := getProductByIDDB(db, 1)
	if prod1 != nil {
		fmt.Printf("2. getProductByIDDB(db, 1): ID: %d | %s | ราคา: %.2f | คงเหลือ: %d\n",
			prod1.ID, prod1.Name, prod1.Price, prod1.Quantity)
	}

	likeProds, _ := getProductByNameDB(db, "Galaxy")
	fmt.Printf("3. getProductByNameDB(db, \"Galaxy\"): พบ %d รายการ:\n", len(likeProds))
	for _, p := range likeProds {
		fmt.Printf("   - %s (%.2f บาท)\n", p.Name, p.Price)
	}

	catProds, _ := getProductsByCategoryDB(db, 1)
	fmt.Printf("4. getProductsByCategoryDB(db, 1): ในหมวดหมู่ Laptops พบ %d รายการ\n", len(catProds))

	rangeProds, _ := getProductsByPriceRangeDB(db, 1000, 10000)
	fmt.Printf("5. getProductsByPriceRangeDB(db, 1000, 10000): พบ %d รายการ ในช่วงราคา 1,000 - 10,000 บาท\n", len(rangeProds))

	lowStockProds, _ := getProductsLowStockDB(db, 5)
	fmt.Printf("6. getProductsLowStockDB(db, 5): สินค้าที่สต็อกน้อยกว่า 5 ชิ้น (พบ %d รายการ):\n", len(lowStockProds))
	for _, p := range lowStockProds {
		fmt.Printf("   - %s คงเหลือ %d ชิ้น\n", p.Name, p.Quantity)
	}

	catPriceProds, _ := getProductsByCategoryAndPriceDB(db, 2, 20000, 35000)
	fmt.Printf("7. getProductsByCategoryAndPriceDB(db, 2, 20000, 35000): พบ %d รายการ\n", len(catPriceProds))

	ascProds, _ := getProductsSortedByPriceDB(db, "asc")
	if len(ascProds) > 0 {
		fmt.Printf("8. getProductsSortedByPriceDB(db, \"asc\"): สินค้าราคาถูกที่สุด คือ '%s' (%.2f บาท)\n",
			ascProds[0].Name, ascProds[0].Price)
	}

	top3, _ := getTopExpensiveProductsDB(db, 3)
	fmt.Printf("9. getTopExpensiveProductsDB(db, 3): สินค้าแพงที่สุด 3 รายการ:\n")
	for i, p := range top3 {
		fmt.Printf("   %d. %s (%.2f บาท)\n", i+1, p.Name, p.Price)
	}

	targetEmail := "somchai.j@example.com"
	custByEmail, _ := getCustomerByEmailDB(db, targetEmail)
	if custByEmail != nil {
		fmt.Printf("10. getCustomerByEmailDB(db, \"%s\"): พบลูกค้า %s (เบอร์โทร: %s)\n",
			targetEmail, custByEmail.Name, custByEmail.Phone)
	}

	custOrders, _ := getOrdersByCustomerDB(db, 1)
	fmt.Printf("11. getOrdersByCustomerDB(db, 1): ลูกค้า ID 1 มีคำสั่งซื้อทั้งหมด %d รายการ\n", len(custOrders))

	pendingOrders, _ := getOrdersByStatusDB(db, "Pending")
	fmt.Printf("12. getOrdersByStatusDB(db, \"Pending\"): มีคำสั่งซื้อสถานะ Pending %d รายการ\n", len(pendingOrders))

	yesterday := time.Now().AddDate(0, 0, -1)
	tomorrow := time.Now().AddDate(0, 0, 1)
	dateOrders, _ := getOrdersByDateRangeDB(db, yesterday, tomorrow)
	fmt.Printf("13. getOrdersByDateRangeDB(db, yesterday, tomorrow): คำสั่งซื้อในช่วงนี้ %d รายการ\n", len(dateOrders))

	detailedOrders, _ := getOrdersWithDetailsDB(db)
	fmt.Printf("14. getOrdersWithDetailsDB(db): ดึงข้อมูลคำสั่งซื้อพร้อมลูกค้าและสินค้า (ทั้งหมด %d รายการ)\n", len(detailedOrders))
	if len(detailedOrders) > 0 {
		o := detailedOrders[0]
		fmt.Printf("    ตัวอย่างรายการแรก: Order #%d | ลูกค้า: %s | สินค้า: %s | ยอดรวม: %.2f | สถานะ: %s\n",
			o.ID, o.Customer.Name, o.Product.Name, o.TotalPrice, o.Status)
	}

	printSubHeader("ส่วนที่ 4: Update (ทดสอบการแก้ไขข้อมูลและเงื่อนไข)")

	updatedProd, err := updateProductByIDDB(db, newProd.ID, "Xiaomi Smart Camera C400 Pro", 1690.00, 15, newCat.ID)
	if err == nil {
		fmt.Printf("1. updateProductByIDDB(): แก้ไขชื่อสินค้าเป็น '%s' ราคา %.2f สต็อก %d\n",
			updatedProd.Name, updatedProd.Price, updatedProd.Quantity)
	}

	_, errPriceNeg := updateProductPriceDB(db, newProd.ID, -500)
	if errPriceNeg != nil {
		fmt.Printf("2. ทดสอบราคาติดลบ: ป้องกันสำเร็จตามเงื่อนไข -> %v\n", errPriceNeg)
	}
	pPriceUpdated, _ := updateProductPriceDB(db, newProd.ID, 1590.00)
	if pPriceUpdated != nil {
		fmt.Printf("   updateProductPriceDB(): อัปเดตราคาใหม่เป็น %.2f บาท\n", pPriceUpdated.Price)
	}

	_, errStockNeg := updateProductStockDB(db, newProd.ID, -999)
	if errStockNeg != nil {
		fmt.Printf("3. ทดสอบสต็อกติดลบ: ป้องกันสำเร็จตามเงื่อนไข -> %v\n", errStockNeg)
	}
	pStockUpdated, _ := updateProductStockDB(db, newProd.ID, +5)
	if pStockUpdated != nil {
		fmt.Printf("   updateProductStockDB(): เพิ่มสต็อก 5 ชิ้น (คงเหลือรวม: %d ชิ้น)\n", pStockUpdated.Quantity)
	}

	updatedEmail := fmt.Sprintf("test.updated.%d@example.com", time.Now().UnixNano())
	custUpdated, _ := updateCustomerByIDDB(db, newCust.ID, "ทดสอบ ผู้ใช้ (อัปเดต)", updatedEmail, "089-123-4567")
	if custUpdated != nil {
		fmt.Printf("4. updateCustomerByIDDB(): อัปเดตข้อมูลลูกค้าเป็น Name: %s, Email: %s\n", custUpdated.Name, custUpdated.Email)
	}

	_, errInvalidStatus := updateOrderStatusDB(db, newOrder.ID, "Shipped")
	if errInvalidStatus != nil {
		fmt.Printf("5. ทดสอบสถานะไม่ถูกต้อง: ป้องกันสำเร็จตามเงื่อนไข -> %v\n", errInvalidStatus)
	}
	orderStatusUpdated, _ := updateOrderStatusDB(db, newOrder.ID, "Completed")
	if orderStatusUpdated != nil {
		fmt.Printf("   updateOrderStatusDB(): เปลี่ยนสถานะคำสั่งซื้อ #%d เป็น '%s' สำเร็จ\n",
			orderStatusUpdated.ID, orderStatusUpdated.Status)
	}

	printSubHeader("ส่วนที่ 5: Delete (ทดสอบ Soft Delete, Hard Delete, Restore และความสัมพันธ์)")

	tempProd, _ := createProductDB(db, "สินค้าทดสอบสำหรับลบ", 500.00, 10, newCat.ID)

	err = softDeleteProductDB(db, tempProd.ID)
	if err == nil {
		fmt.Printf("1. softDeleteProductDB(): Soft Delete สินค้า ID %d สำเร็จ\n", tempProd.ID)
	}

	delProducts, _ := getDeletedProductsDB(db)
	fmt.Printf("2. getDeletedProductsDB(): พบสินค้าที่ถูก Soft Delete ทั้งหมด %d รายการ\n", len(delProducts))

	err = restoreProductDB(db, tempProd.ID)
	if err == nil {
		fmt.Printf("3. restoreProductDB(): กู้คืนสินค้า ID %d สำเร็จ\n", tempProd.ID)
	}

	err = hardDeleteProductDB(db, tempProd.ID)
	if err == nil {
		fmt.Printf("4. hardDeleteProductDB(): ลบสินค้า ID %d ถาวรสำเร็จ (ไม่มีคำสั่งซื้ออ้างอิง)\n", tempProd.ID)
	}

	errProdWithOrder := hardDeleteProductDB(db, newProd.ID)
	if errProdWithOrder != nil {
		fmt.Printf("5. ทดสอบ Hard Delete สินค้าที่มี Order: ป้องกันสำเร็จ -> %v\n", errProdWithOrder)
	}

	errCatWithProd := deleteCategoryByIDDB(db, newCat.ID)
	if errCatWithProd != nil {
		fmt.Printf("6. ทดสอบลบหมวดหมู่ที่ยังมีสินค้า: ป้องกันสำเร็จ -> %v\n", errCatWithProd)
	}

	errCustWithOrder := deleteCustomerByIDDB(db, newCust.ID)
	if errCustWithOrder != nil {
		fmt.Printf("7. ทดสอบลบลูกค้าที่มี Order: ป้องกันสำเร็จ -> %v\n", errCustWithOrder)
	}

	printDivider("การทดสอบทุกฟังก์ชันเสร็จสมบูรณ์เรียบร้อยแล้ว!")
}
