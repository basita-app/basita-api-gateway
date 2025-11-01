package cms

// GraphQL queries for all CMS operations

const (
	// GetBrandsQuery fetches all brands with simplified response
	GetBrandsQuery = `
		query GetBrands {
			brands {
				documentId
				Name
				Logo {
					documentId
					url
					width
					height
					formats
				}
			}
		}
	`

	// GetBrandByIDQuery fetches a single brand by documentId
	GetBrandByIDQuery = `
		query GetBrand($documentId: ID!) {
			brand(documentId: $documentId) {
				documentId
				Name
				Slug
				Logo {
					documentId
					url
					width
					height
					formats
				}
			}
		}
	`

	// GetCarModelsByBrandQuery fetches car models for a specific brand
	GetCarModelsByBrandQuery = `
		query GetCarModelsByBrand($brandDocumentId: ID!) {
			carModels(filters: { brand: { documentId: { eq: $brandDocumentId } } }) {
				documentId
				Name
				Images {
					documentId
					url
					width
					height
					formats
				}
			}
		}
	`

	// GetCarVariantsByModelQuery fetches variants with showrooms for a car model
	GetCarVariantsByModelQuery = `
		query GetCarVariants($carModelDocumentId: ID!) {
			carVariants(filters: { car_model: { documentId: { eq: $carModelDocumentId } } }) {
				documentId
				Name
				Price
				Year
				BrochureURL
				ReviewLink
				Warranty
				MinimumDownPaymet
				MinimumInstallments
				Specs {
					Motor
				}
				ShowroomPricing {
					Price
					showroom {
						documentId
						Name
						Logo {
							documentId
							url
							width
							height
							formats
						}
					}
				}
			}
		}
	`

	// GetCarModelByIDQuery fetches a single car model with images
	GetCarModelByIDQuery = `
		query GetCarModel($documentId: ID!) {
			carModel(documentId: $documentId) {
				documentId
				Name
				BodyType
				FuelType
				Slug
				Images {
					documentId
					url
					width
					height
					formats
				}
			}
		}
	`

	// GetAdvertisementsQuery fetches all advertisements
	GetAdvertisementsQuery = `
		query GetAdvertisements {
			advertisements {
				documentId
				Action
				Banner {
					url
					documentId
					name
					width
					height
					formats
				}
			}
		}
	`

	// GetAdvertisementByIDQuery fetches a single advertisement
	GetAdvertisementByIDQuery = `
		query GetAdvertisement($documentId: ID!) {
			advertisement(documentId: $documentId) {
				documentId
				Action
				Banner {
					url
					documentId
					name
					width
					height
					formats
				}
			}
		}
	`

	// GetCarVariantByIDQuery fetches a single car variant with all details
	GetCarVariantByIDQuery = `
		query GetCarVariant($documentId: ID!) {
			carVariant(documentId: $documentId) {
				documentId
				Name
				Price
				Year
				BrochureURL
				ReviewLink
				Warranty
				MinimumDownPaymet
				MinimumInstallments
				car_model {
					documentId
					Name
					Images {
						documentId
						url
						width
						height
						formats
					}
				}
				Specs {
					Motor
					Transmission
					Acceleration
					AssembledIn
					GroundClearanceInMM
					HeightInMM
					Horsepower
					LengthInMM
					LiterPerKM
					MaxSpeed
					Origin
					Speed
					TractionType
					Transmission
					TrunkSize
					WheelBase
					WidthInMM
					Seats
				}
				Features
				ShowroomPricing {
					Price
					MinimuDownpayment
					MinimumInstallements
					showroom {
						documentId
						Name
						Logo {
							documentId
							url
							width
							height
							formats
						}
					}
				}
			}
		}
	`

	// GetShowroomsQuery fetches all showrooms
	GetShowroomsQuery = `
		query GetShowrooms {
			showrooms {
				documentId
				Name
				Description
				IsVerified
				IsFeatured
				Logo {
					documentId
					url
					width
					height
					formats
				}
			}
		}
	`

	// GetShowroomByIDQuery fetches a single showroom by documentId with full details
	GetShowroomByIDQuery = `
		query GetShowroomProfile($documentId: ID!) {
			showroom(documentId: $documentId) {
				documentId
				Logo {
					documentId
					url
					width
					height
					formats
				}
				Cover {
					documentId
					url
					width
					height
					formats
				}
				Name
				Description
				IsVerified
				IsFeatured
				OperatingHours
				Location {
					Address
					governorate {
						documentId
						Name
					}
					city {
						documentId
						Name
					}
					Latitude
					Longitude
				}
				ContactInfo {
					Email
					Phone
					Facebook
					Instagram
					Tiktok
					Whatsapp
					X
					Youtube
					WebsiteURL
				}
			}
		}
	`

	// GetCarVariantsByShowroomQuery fetches car variants for a specific showroom
	GetCarVariantsByShowroomQuery = `
		query GetCarVariantsByShowroom($showroomDocumentId: ID!) {
			carVariants(filters: { ShowroomPricing: { showroom: { documentId: { eq: $showroomDocumentId } } } }) {
				documentId
				DisplayName
				Images {
					documentId
						url
						width
						height
						formats
				}
				ShowroomPricing(filters: { showroom: { documentId: { eq: $showroomDocumentId } } }) {
					Price
					MinimuDownpayment
					MinimumInstallements
				}
			}
		}`
)
