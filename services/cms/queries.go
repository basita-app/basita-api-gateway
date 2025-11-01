package cms

// GraphQL queries for all CMS operations

const (
	// GetBrandsQuery fetches all brands with simplified response
	GetBrandsQuery = `
		query GetBrands($locale: I18NLocaleCode) {
			brands(locale: $locale) {
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
		query GetBrand($documentId: ID!, $locale: I18NLocaleCode) {
			brand(documentId: $documentId, locale: $locale) {
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
		query GetCarModelsByBrand($brandDocumentId: ID!, $locale: I18NLocaleCode) {
			carModels(filters: { brand: { documentId: { eq: $brandDocumentId } } }, locale: $locale) {
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
		query GetCarVariants($carModelDocumentId: ID!, $locale: I18NLocaleCode) {
			carVariants(filters: { car_model: { documentId: { eq: $carModelDocumentId } } }, locale: $locale) {
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
		query GetCarModel($documentId: ID!, $locale: I18NLocaleCode) {
			carModel(documentId: $documentId, locale: $locale) {
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
		query GetAdvertisements($locale: I18NLocaleCode) {
			advertisements(locale: $locale) {
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
		query GetAdvertisement($documentId: ID!, $locale: I18NLocaleCode) {
			advertisement(documentId: $documentId, locale: $locale) {
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
		query GetCarVariant($documentId: ID!, $locale: I18NLocaleCode) {
			carVariant(documentId: $documentId, locale: $locale) {
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
		query GetShowrooms($locale: I18NLocaleCode) {
			showrooms(locale: $locale) {
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
		query GetShowroomProfile($documentId: ID!, $locale: I18NLocaleCode) {
			showroom(documentId: $documentId, locale: $locale) {
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
		query GetCarVariantsByShowroom($showroomDocumentId: ID!, $locale: I18NLocaleCode) {
			carVariants(filters: { ShowroomPricing: { showroom: { documentId: { eq: $showroomDocumentId } } } }, locale: $locale) {
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
