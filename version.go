package main

const (
	version string = "1.2.3"
	// 05/06/2020
	// Update aldo product with mongodbId if zunka product alredy pointing to aldo product.

	// version string = "1.2.2"
	// // 05/06/2020
	// // Bug fix - Not calling check consistency routine.

	// version string = "1.2.1"
	// // 04/06/2020
	// // Check consistency for products aldo and zunkasite products.

	// version string = "1.2.0"
	// // Update price and availability for products created at zunkasite.

	// version string = "1.1.0"
	// // Using new db version without fields "new", "removed" and "changed".

	// version string = "1.0.1"
	// // Bugfix - When product changed mondodbId was not copied from old product.
	// // Bugfix - Removed _new from table names scripts.

	// // Not using tables product and product history without id
	// version string = "1.0.0"

	// // Better log for sqlx.Exec().
	// version string = "0.5.3"

	// Min price 1050,00.
	// version string = "0.5.1"
)
