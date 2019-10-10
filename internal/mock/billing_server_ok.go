package mock

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/micro/go-micro/client"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
)

type BillingServerOkMock struct{}

func NewBillingServerOkMock() grpc.BillingService {
	return &BillingServerOkMock{}
}

func (s *BillingServerOkMock) GetProductsForOrder(
	ctx context.Context,
	in *grpc.GetProductsForOrderRequest,
	opts ...client.CallOption,
) (*grpc.ListProductsResponse, error) {
	return &grpc.ListProductsResponse{}, nil
}

func (s *BillingServerOkMock) OrderCreateProcess(
	ctx context.Context,
	in *billing.OrderCreateRequest,
	opts ...client.CallOption,
) (*grpc.OrderCreateProcessResponse, error) {
	return &grpc.OrderCreateProcessResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Order{},
	}, nil
}

func (s *BillingServerOkMock) PaymentFormJsonDataProcess(
	ctx context.Context,
	in *grpc.PaymentFormJsonDataRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormJsonDataResponse, error) {
	cookie := in.Cookie
	if cookie == "" {
		cookie = bson.NewObjectId().Hex()
	}
	return &grpc.PaymentFormJsonDataResponse{
		Status: pkg.ResponseStatusOk,
		Item: &grpc.PaymentFormJsonData{
			Cookie: cookie,
		},
	}, nil
}

func (s *BillingServerOkMock) PaymentCreateProcess(
	ctx context.Context,
	in *grpc.PaymentCreateRequest,
	opts ...client.CallOption,
) (*grpc.PaymentCreateResponse, error) {
	return &grpc.PaymentCreateResponse{}, nil
}

func (s *BillingServerOkMock) PaymentCallbackProcess(
	ctx context.Context,
	in *grpc.PaymentNotifyRequest,
	opts ...client.CallOption,
) (*grpc.PaymentNotifyResponse, error) {
	return &grpc.PaymentNotifyResponse{}, nil
}

func (s *BillingServerOkMock) RebuildCache(
	ctx context.Context,
	in *grpc.EmptyRequest,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkMock) UpdateOrder(
	ctx context.Context,
	in *billing.Order,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkMock) UpdateMerchant(
	ctx context.Context,
	in *billing.Merchant,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkMock) GetMerchantBy(
	ctx context.Context,
	in *grpc.GetMerchantByRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantResponse, error) {
	if in.MerchantId == OnboardingMerchantMock.Id {
		OnboardingMerchantMock.S3AgreementName = SomeAgreementName
	} else {
		if in.MerchantId == SomeMerchantId1 {
			OnboardingMerchantMock.S3AgreementName = SomeAgreementName1
		} else {
			if in.MerchantId == SomeMerchantId2 {
				OnboardingMerchantMock.S3AgreementName = SomeAgreementName2
			} else {
				OnboardingMerchantMock.S3AgreementName = ""
			}
		}
	}

	if in.MerchantId == SomeMerchantId3 {
		OnboardingMerchantMock.Status = pkg.MerchantStatusDraft
	} else {
		OnboardingMerchantMock.Status = pkg.MerchantStatusAgreementSigning
	}

	rsp := &grpc.GetMerchantResponse{
		Status:  pkg.ResponseStatusOk,
		Message: &grpc.ResponseErrorMessage{},
		Item:    OnboardingMerchantMock,
	}

	return rsp, nil
}

func (s *BillingServerOkMock) ListMerchants(
	ctx context.Context,
	in *grpc.MerchantListingRequest,
	opts ...client.CallOption,
) (*grpc.MerchantListingResponse, error) {
	return &grpc.MerchantListingResponse{
		Count: 3,
		Items: []*billing.Merchant{OnboardingMerchantMock, OnboardingMerchantMock, OnboardingMerchantMock},
	}, nil
}

func (s *BillingServerOkMock) ChangeMerchant(
	ctx context.Context,
	in *grpc.OnboardingRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantResponse, error) {
	m := &billing.Merchant{
		User: &billing.MerchantUser{
			Id:    bson.NewObjectId().Hex(),
			Email: "test@unit.test",
		},
		Company:  in.Company,
		Contacts: in.Contacts,
		Banking:  in.Banking,
		Status:   pkg.MerchantStatusDraft,
	}

	if in.Id != "" {
		m.Id = in.Id
	} else {
		m.Id = bson.NewObjectId().Hex()
	}

	return &grpc.ChangeMerchantResponse{
		Status: pkg.ResponseStatusOk,
		Item:   m,
	}, nil
}

func (s *BillingServerOkMock) ChangeMerchantStatus(
	ctx context.Context,
	in *grpc.MerchantChangeStatusRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantStatusResponse, error) {
	return &grpc.ChangeMerchantStatusResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Merchant{Id: in.MerchantId, Status: in.Status},
	}, nil
}

func (s *BillingServerOkMock) CreateNotification(
	ctx context.Context,
	in *grpc.NotificationRequest,
	opts ...client.CallOption,
) (*grpc.CreateNotificationResponse, error) {
	return &grpc.CreateNotificationResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) GetNotification(
	ctx context.Context,
	in *grpc.GetNotificationRequest,
	opts ...client.CallOption,
) (*billing.Notification, error) {
	return &billing.Notification{}, nil
}

func (s *BillingServerOkMock) ListNotifications(
	ctx context.Context,
	in *grpc.ListingNotificationRequest,
	opts ...client.CallOption,
) (*grpc.Notifications, error) {
	return &grpc.Notifications{}, nil
}

func (s *BillingServerOkMock) MarkNotificationAsRead(
	ctx context.Context,
	in *grpc.GetNotificationRequest,
	opts ...client.CallOption,
) (*billing.Notification, error) {
	return &billing.Notification{}, nil
}

func (s *BillingServerOkMock) ListMerchantPaymentMethods(
	ctx context.Context,
	in *grpc.ListMerchantPaymentMethodsRequest,
	opts ...client.CallOption,
) (*grpc.ListingMerchantPaymentMethod, error) {
	return &grpc.ListingMerchantPaymentMethod{}, nil
}

func (s *BillingServerOkMock) GetMerchantPaymentMethod(
	ctx context.Context,
	in *grpc.GetMerchantPaymentMethodRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantPaymentMethodResponse, error) {
	return &grpc.GetMerchantPaymentMethodResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) ChangeMerchantPaymentMethod(
	ctx context.Context,
	in *grpc.MerchantPaymentMethodRequest,
	opts ...client.CallOption,
) (*grpc.MerchantPaymentMethodResponse, error) {
	return &grpc.MerchantPaymentMethodResponse{
		Status: pkg.ResponseStatusOk,
		Item: &billing.MerchantPaymentMethod{
			PaymentMethod: &billing.MerchantPaymentMethodIdentification{
				Id:   in.PaymentMethod.Id,
				Name: in.PaymentMethod.Name,
			},
			Commission:  in.Commission,
			Integration: in.Integration,
			IsActive:    in.IsActive,
		},
	}, nil
}

func (s *BillingServerOkMock) CreateRefund(
	ctx context.Context,
	in *grpc.CreateRefundRequest,
	opts ...client.CallOption,
) (*grpc.CreateRefundResponse, error) {
	return &grpc.CreateRefundResponse{
		Status: pkg.ResponseStatusOk,
		Item: &billing.Refund{
			Id:            bson.NewObjectId().Hex(),
			OriginalOrder: &billing.RefundOrder{Id: bson.NewObjectId().Hex(), Uuid: uuid.New().String()},
			ExternalId:    "",
			Amount:        10,
			CreatorId:     "",
			Reason:        SomeError.Message,
			Currency:      "RUB",
			Status:        0,
		},
	}, nil
}

func (s *BillingServerOkMock) ListRefunds(
	ctx context.Context,
	in *grpc.ListRefundsRequest,
	opts ...client.CallOption,
) (*grpc.ListRefundsResponse, error) {
	return &grpc.ListRefundsResponse{
		Count: 2,
		Items: []*billing.Refund{
			{
				Id:            bson.NewObjectId().Hex(),
				OriginalOrder: &billing.RefundOrder{Id: bson.NewObjectId().Hex(), Uuid: uuid.New().String()},
				ExternalId:    "",
				Amount:        10,
				CreatorId:     "",
				Reason:        SomeError.Message,
				Currency:      "RUB",
				Status:        0,
			},
			{
				Id:            bson.NewObjectId().Hex(),
				OriginalOrder: &billing.RefundOrder{Id: bson.NewObjectId().Hex(), Uuid: uuid.New().String()},
				ExternalId:    "",
				Amount:        10,
				CreatorId:     "",
				Reason:        SomeError.Message,
				Currency:      "RUB",
				Status:        0,
			},
		},
	}, nil
}

func (s *BillingServerOkMock) GetRefund(
	ctx context.Context,
	in *grpc.GetRefundRequest,
	opts ...client.CallOption,
) (*grpc.CreateRefundResponse, error) {
	return &grpc.CreateRefundResponse{
		Status: pkg.ResponseStatusOk,
		Item: &billing.Refund{
			Id:            bson.NewObjectId().Hex(),
			OriginalOrder: &billing.RefundOrder{Id: bson.NewObjectId().Hex(), Uuid: uuid.New().String()},
			ExternalId:    "",
			Amount:        10,
			CreatorId:     "",
			Reason:        SomeError.Message,
			Currency:      "RUB",
			Status:        0,
		},
	}, nil
}

func (s *BillingServerOkMock) ProcessRefundCallback(
	ctx context.Context,
	in *grpc.CallbackRequest,
	opts ...client.CallOption,
) (*grpc.PaymentNotifyResponse, error) {
	return &grpc.PaymentNotifyResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) PaymentFormLanguageChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangeLangRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormDataChangeResponse, error) {
	return &grpc.PaymentFormDataChangeResponse{
		Status: pkg.ResponseStatusOk,
		Item: &grpc.PaymentFormDataChangeResponseItem{
			UserAddressDataRequired: true,
			UserIpData: &grpc.UserIpData{
				Country: "RU",
				City:    "St.Petersburg",
				Zip:     "190000",
			},
		},
	}, nil
}

func (s *BillingServerOkMock) PaymentFormPaymentAccountChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangePaymentAccountRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormDataChangeResponse, error) {
	return &grpc.PaymentFormDataChangeResponse{
		Status: pkg.ResponseStatusOk,
		Item: &grpc.PaymentFormDataChangeResponseItem{
			UserAddressDataRequired: true,
			UserIpData: &grpc.UserIpData{
				Country: "RU",
				City:    "St.Petersburg",
				Zip:     "190000",
			},
		},
	}, nil
}

func (s *BillingServerOkMock) ProcessBillingAddress(
	ctx context.Context,
	in *grpc.ProcessBillingAddressRequest,
	opts ...client.CallOption,
) (*grpc.ProcessBillingAddressResponse, error) {
	return &grpc.ProcessBillingAddressResponse{
		Status: pkg.ResponseStatusOk,
		Item: &grpc.ProcessBillingAddressResponseItem{
			HasVat:      true,
			Vat:         10,
			Amount:      10,
			TotalAmount: 20,
		},
	}, nil
}

func (s *BillingServerOkMock) ChangeMerchantData(
	ctx context.Context,
	in *grpc.ChangeMerchantDataRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantDataResponse, error) {
	rsp := &grpc.ChangeMerchantDataResponse{
		Status: pkg.ResponseStatusOk,
		Item:   OnboardingMerchantMock,
	}

	if in.MerchantId == SomeMerchantId {
		return nil, SomeError
	}

	return rsp, nil
}

func (s *BillingServerOkMock) SetMerchantS3Agreement(
	ctx context.Context,
	in *grpc.SetMerchantS3AgreementRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantDataResponse, error) {
	rsp := &grpc.ChangeMerchantDataResponse{
		Status: pkg.ResponseStatusOk,
		Item:   OnboardingMerchantMock,
	}

	if in.MerchantId == SomeMerchantId {
		return nil, SomeError
	}

	return rsp, nil
}

func (s *BillingServerOkMock) ChangeProject(
	ctx context.Context,
	in *billing.Project,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return &grpc.ChangeProjectResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetProject(
	ctx context.Context,
	in *grpc.GetProjectRequest,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return &grpc.ChangeProjectResponse{
		Status: pkg.ResponseStatusOk,
		Item: &billing.Project{
			MerchantId:         bson.NewObjectId().Hex(),
			Name:               map[string]string{"en": "A", "ru": "А"},
			CallbackCurrency:   "RUB",
			CallbackProtocol:   pkg.ProjectCallbackProtocolEmpty,
			LimitsCurrency:     "RUB",
			MinPaymentAmount:   0,
			MaxPaymentAmount:   15000,
			IsProductsCheckout: false,
		},
	}, nil
}

func (s *BillingServerOkMock) DeleteProject(
	ctx context.Context,
	in *grpc.GetProjectRequest,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return &grpc.ChangeProjectResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) CreateToken(
	ctx context.Context,
	in *grpc.TokenRequest,
	opts ...client.CallOption,
) (*grpc.TokenResponse, error) {
	return &grpc.TokenResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) CheckProjectRequestSignature(
	ctx context.Context,
	in *grpc.CheckProjectRequestSignatureRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) CreateOrUpdateProduct(ctx context.Context, in *grpc.Product, opts ...client.CallOption) (*grpc.Product, error) {
	return Product, nil
}

func (s *BillingServerOkMock) ListProducts(ctx context.Context, in *grpc.ListProductsRequest, opts ...client.CallOption) (*grpc.ListProductsResponse, error) {
	return &grpc.ListProductsResponse{
		Limit:  1,
		Offset: 0,
		Total:  200,
		Products: []*grpc.Product{
			Product,
		},
	}, nil
}

func (s *BillingServerOkMock) GetProduct(ctx context.Context, in *grpc.RequestProduct, opts ...client.CallOption) (*grpc.GetProductResponse, error) {
	return GetProductResponse, nil
}

func (s *BillingServerOkMock) DeleteProduct(ctx context.Context, in *grpc.RequestProduct, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkMock) FindAllOrders(ctx context.Context, in *grpc.ListOrdersRequest, opts ...client.CallOption) (*billing.OrderPaginate, error) {
	return &billing.OrderPaginate{Count: 1, Items: []*billing.Order{}}, nil
}

func (s *BillingServerOkMock) ListProjects(ctx context.Context, in *grpc.ListProjectsRequest, opts ...client.CallOption) (*grpc.ListProjectsResponse, error) {
	return &grpc.ListProjectsResponse{Count: 1, Items: []*billing.Project{{Id: "id"}}}, nil
}
func (s *BillingServerOkMock) GetOrder(ctx context.Context, in *grpc.GetOrderRequest, opts ...client.CallOption) (*billing.Order, error) {
	return &billing.Order{}, nil
}

func (s *BillingServerOkMock) IsOrderCanBePaying(
	ctx context.Context,
	in *grpc.IsOrderCanBePayingRequest,
	opts ...client.CallOption,
) (*grpc.IsOrderCanBePayingResponse, error) {
	return &grpc.IsOrderCanBePayingResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Order{},
	}, nil
}

func (s *BillingServerOkMock) GetCountry(ctx context.Context, in *billing.GetCountryRequest, opts ...client.CallOption) (*billing.Country, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) UpdateCountry(ctx context.Context, in *billing.Country, opts ...client.CallOption) (*billing.Country, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetPriceGroup(ctx context.Context, in *billing.GetPriceGroupRequest, opts ...client.CallOption) (*billing.PriceGroup, error) {
	return &billing.PriceGroup{
		Id: "some_id",
	}, nil
}

func (s *BillingServerOkMock) UpdatePriceGroup(ctx context.Context, in *billing.PriceGroup, opts ...client.CallOption) (*billing.PriceGroup, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) SetUserNotifySales(ctx context.Context, in *grpc.SetUserNotifyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) SetUserNotifyNewRegion(ctx context.Context, in *grpc.SetUserNotifyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetCountriesList(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*billing.CountriesList, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetPaymentChannelCostSystem(ctx context.Context, in *billing.PaymentChannelCostSystemRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemResponse, error) {
	return &grpc.PaymentChannelCostSystemResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) SetPaymentChannelCostSystem(ctx context.Context, in *billing.PaymentChannelCostSystem, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemResponse, error) {
	return &grpc.PaymentChannelCostSystemResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) DeletePaymentChannelCostSystem(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchantRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantResponse, error) {
	return &grpc.PaymentChannelCostMerchantResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) SetPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchant, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantResponse, error) {
	return &grpc.PaymentChannelCostMerchantResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) DeletePaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetMoneyBackCostSystem(ctx context.Context, in *billing.MoneyBackCostSystemRequest, opts ...client.CallOption) (*grpc.MoneyBackCostSystemResponse, error) {
	return &grpc.MoneyBackCostSystemResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) SetMoneyBackCostSystem(ctx context.Context, in *billing.MoneyBackCostSystem, opts ...client.CallOption) (*grpc.MoneyBackCostSystemResponse, error) {
	return &grpc.MoneyBackCostSystemResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) DeleteMoneyBackCostSystem(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchantRequest, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantResponse, error) {
	return &grpc.MoneyBackCostMerchantResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.MoneyBackCostMerchant{},
	}, nil
}

func (s *BillingServerOkMock) SetMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchant, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantResponse, error) {
	return &grpc.MoneyBackCostMerchantResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.MoneyBackCostMerchant{},
	}, nil
}

func (s *BillingServerOkMock) DeleteMoneyBackCostMerchant(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) CreateOrUpdatePaymentMethodTestSettings(ctx context.Context, in *grpc.ChangePaymentMethodParamsRequest, opts ...client.CallOption) (*grpc.ChangePaymentMethodParamsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) DeletePaymentMethodTestSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.ChangePaymentMethodParamsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) FindByZipCode(
	ctx context.Context,
	in *grpc.FindByZipCodeRequest,
	opts ...client.CallOption,
) (*grpc.FindByZipCodeResponse, error) {
	return &grpc.FindByZipCodeResponse{
		Count: 1,
		Items: []*billing.ZipCode{
			{
				Zip:     in.Zip,
				Country: in.Country,
			},
		},
	}, nil
}

func (s *BillingServerOkMock) GetAllPaymentChannelCostSystem(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemListResponse, error) {
	return &grpc.PaymentChannelCostSystemListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) GetAllPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchantListRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantListResponse, error) {
	return &grpc.PaymentChannelCostMerchantListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) GetAllMoneyBackCostSystem(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.MoneyBackCostSystemListResponse, error) {
	return &grpc.MoneyBackCostSystemListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) GetAllMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchantListRequest, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantListResponse, error) {
	return &grpc.MoneyBackCostMerchantListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) CreateOrUpdatePaymentMethod(
	ctx context.Context,
	in *billing.PaymentMethod,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodResponse, error) {
	return &grpc.ChangePaymentMethodResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) CreateOrUpdatePaymentMethodProductionSettings(
	ctx context.Context,
	in *grpc.ChangePaymentMethodParamsRequest,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodParamsResponse, error) {
	return &grpc.ChangePaymentMethodParamsResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) DeletePaymentMethodProductionSettings(
	ctx context.Context,
	in *grpc.GetPaymentMethodSettingsRequest,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodParamsResponse, error) {
	return &grpc.ChangePaymentMethodParamsResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) CreateAccountingEntry(ctx context.Context, in *grpc.CreateAccountingEntryRequest, opts ...client.CallOption) (*grpc.CreateAccountingEntryResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) CreateRoyaltyReport(ctx context.Context, in *grpc.CreateRoyaltyReportRequest, opts ...client.CallOption) (*grpc.CreateRoyaltyReportRequest, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) ListRoyaltyReports(ctx context.Context, in *grpc.ListRoyaltyReportsRequest, opts ...client.CallOption) (*grpc.ListRoyaltyReportsResponse, error) {
	return &grpc.ListRoyaltyReportsResponse{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
		Data: &grpc.RoyaltyReportsPaginate{
			Count: 1,
			Items: []*billing.RoyaltyReport{},
		},
	}, nil
}

func (s *BillingServerOkMock) ListRoyaltyReportOrders(ctx context.Context, in *grpc.ListRoyaltyReportOrdersRequest, opts ...client.CallOption) (*grpc.TransactionsResponse, error) {
	return &grpc.TransactionsResponse{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
		Data: &grpc.TransactionsPaginate{
			Count: 100,
			Items: []*billing.OrderViewPublic{},
		},
	}, nil
}

func (s *BillingServerOkMock) GetVatReportsDashboard(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.VatReportsResponse, error) {
	return &grpc.VatReportsResponse{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
		Data: &grpc.VatReportsPaginate{
			Count: 1,
			Items: []*billing.VatReport{},
		},
	}, nil
}

func (s *BillingServerOkMock) GetVatReportsForCountry(ctx context.Context, in *grpc.VatReportsRequest, opts ...client.CallOption) (*grpc.VatReportsResponse, error) {
	return &grpc.VatReportsResponse{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
		Data: &grpc.VatReportsPaginate{
			Count: 100,
			Items: []*billing.VatReport{},
		},
	}, nil
}

func (s *BillingServerOkMock) GetVatReportTransactions(ctx context.Context, in *grpc.VatTransactionsRequest, opts ...client.CallOption) (*grpc.TransactionsResponse, error) {
	return &grpc.TransactionsResponse{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
		Data: &grpc.TransactionsPaginate{
			Count: 100,
			Items: []*billing.OrderViewPublic{},
		},
	}, nil
}

func (s *BillingServerOkMock) CalcAnnualTurnovers(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}
func (s *BillingServerOkMock) ProcessVatReports(ctx context.Context, in *grpc.ProcessVatReportsRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) UpdateVatReportStatus(ctx context.Context, in *grpc.UpdateVatReportStatusRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
	}, nil
}

func (s *BillingServerOkMock) GetPriceGroupByCountry(
	ctx context.Context,
	in *grpc.PriceGroupByCountryRequest,
	opts ...client.CallOption,
) (*billing.PriceGroup, error) {
	return &billing.PriceGroup{}, nil
}

func (s *BillingServerOkMock) UpdateProductPrices(
	ctx context.Context,
	in *grpc.UpdateProductPricesRequest,
	opts ...client.CallOption,
) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{}, nil
}

func (s *BillingServerOkMock) GetProductPrices(
	ctx context.Context,
	in *grpc.RequestProduct,
	opts ...client.CallOption,
) (*grpc.ProductPricesResponse, error) {
	return &grpc.ProductPricesResponse{}, nil
}

func (s *BillingServerOkMock) GetPriceGroupRecommendedPrice(
	ctx context.Context,
	in *grpc.RecommendedPriceRequest,
	opts ...client.CallOption,
) (*grpc.RecommendedPriceResponse, error) {
	return &grpc.RecommendedPriceResponse{}, nil
}

func (s *BillingServerOkMock) GetPriceGroupCurrencyByRegion(
	ctx context.Context,
	in *grpc.PriceGroupByRegionRequest,
	opts ...client.CallOption,
) (*grpc.PriceGroupCurrenciesResponse, error) {
	return &grpc.PriceGroupCurrenciesResponse{
		Region: []*grpc.PriceGroupRegions{
			{Currency: "USD"},
		},
	}, nil
}

func (s *BillingServerOkMock) GetPriceGroupCurrencies(
	ctx context.Context,
	in *grpc.EmptyRequest,
	opts ...client.CallOption,
) (*grpc.PriceGroupCurrenciesResponse, error) {
	return &grpc.PriceGroupCurrenciesResponse{}, nil
}

func (s *BillingServerOkMock) GetPaymentMethodProductionSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.GetPaymentMethodSettingsResponse, error) {
	return &grpc.GetPaymentMethodSettingsResponse{
		Params: []*billing.PaymentMethodParams{
			{Currency: "RUB"},
		},
	}, nil
}

func (s *BillingServerOkMock) GetPaymentMethodTestSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.GetPaymentMethodSettingsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) ChangeRoyaltyReport(ctx context.Context, in *grpc.ChangeRoyaltyReportRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) AutoAcceptRoyaltyReports(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetUserProfile(
	ctx context.Context,
	in *grpc.GetUserProfileRequest,
	opts ...client.CallOption,
) (*grpc.GetUserProfileResponse, error) {
	return &grpc.GetUserProfileResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &grpc.UserProfile{},
	}, nil
}

func (s *BillingServerOkMock) CreateOrUpdateUserProfile(
	ctx context.Context,
	in *grpc.UserProfile,
	opts ...client.CallOption,
) (*grpc.GetUserProfileResponse, error) {
	return &grpc.GetUserProfileResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &grpc.UserProfile{},
	}, nil
}

func (s *BillingServerOkMock) ConfirmUserEmail(
	ctx context.Context,
	in *grpc.ConfirmUserEmailRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) CreatePageReview(
	ctx context.Context,
	in *grpc.CreatePageReviewRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) MerchantReviewRoyaltyReport(ctx context.Context, in *grpc.MerchantReviewRoyaltyReportRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	return &grpc.ResponseError{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetMerchantAgreementSignUrl(
	ctx context.Context,
	in *grpc.GetMerchantAgreementSignUrlRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantAgreementSignUrlResponse, error) {
	return &grpc.GetMerchantAgreementSignUrlResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetMerchantOnboardingCompleteData(
	ctx context.Context,
	in *grpc.SetMerchantS3AgreementRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantOnboardingCompleteDataResponse, error) {
	return &grpc.GetMerchantOnboardingCompleteDataResponse{Status: pkg.ResponseStatusOk}, nil
}

func (s *BillingServerOkMock) GetMerchantTariffRates(
	ctx context.Context,
	in *grpc.GetMerchantTariffRatesRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantTariffRatesResponse, error) {
	return &grpc.GetMerchantTariffRatesResponse{}, nil
}

func (s *BillingServerOkMock) SetMerchantTariffRates(
	ctx context.Context,
	in *grpc.SetMerchantTariffRatesRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{}, nil
}

func (s *BillingServerOkMock) CreateOrUpdateKeyProduct(ctx context.Context, in *grpc.CreateOrUpdateKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	return &grpc.KeyProductResponse{
		Status:  pkg.ResponseStatusOk,
		Product: &grpc.KeyProduct{},
	}, nil
}

func (s *BillingServerOkMock) GetKeyProducts(ctx context.Context, in *grpc.ListKeyProductsRequest, opts ...client.CallOption) (*grpc.ListKeyProductsResponse, error) {
	return &grpc.ListKeyProductsResponse{
		Status: pkg.ResponseStatusOk,
		Count:  1,
		Products: []*grpc.KeyProduct{
			{},
		},
	}, nil
}

func (s *BillingServerOkMock) GetKeyProduct(ctx context.Context, in *grpc.RequestKeyProductMerchant, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	return &grpc.KeyProductResponse{
		Status:  pkg.ResponseStatusOk,
		Product: &grpc.KeyProduct{},
	}, nil
}

func (s *BillingServerOkMock) DeleteKeyProduct(ctx context.Context, in *grpc.RequestKeyProductMerchant, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	return &grpc.EmptyResponseWithStatus{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) PublishKeyProduct(ctx context.Context, in *grpc.PublishKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	return &grpc.KeyProductResponse{
		Status:  pkg.ResponseStatusOk,
		Product: &grpc.KeyProduct{},
	}, nil
}

func (s *BillingServerOkMock) GetKeyProductsForOrder(ctx context.Context, in *grpc.GetKeyProductsForOrderRequest, opts ...client.CallOption) (*grpc.ListKeyProductsResponse, error) {
	return &grpc.ListKeyProductsResponse{
		Status: pkg.ResponseStatusOk,
		Count:  1,
		Products: []*grpc.KeyProduct{
			{},
		},
	}, nil
}

func (s *BillingServerOkMock) GetPlatforms(ctx context.Context, in *grpc.ListPlatformsRequest, opts ...client.CallOption) (*grpc.ListPlatformsResponse, error) {
	return &grpc.ListPlatformsResponse{
		Status: pkg.ResponseStatusOk,
		Count:  1,
		Platforms: []*grpc.Platform{
			{},
		},
	}, nil
}

func (s *BillingServerOkMock) DeletePlatformFromProduct(ctx context.Context, in *grpc.RemovePlatformRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	return &grpc.EmptyResponseWithStatus{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) GetAvailableKeysCount(ctx context.Context, in *grpc.GetPlatformKeyCountRequest, opts ...client.CallOption) (*grpc.GetPlatformKeyCountResponse, error) {
	return &grpc.GetPlatformKeyCountResponse{
		Status: pkg.ResponseStatusOk,
		Count:  1000,
	}, nil
}

func (s *BillingServerOkMock) UploadKeysFile(ctx context.Context, in *grpc.PlatformKeysFileRequest, opts ...client.CallOption) (*grpc.PlatformKeysFileResponse, error) {
	return &grpc.PlatformKeysFileResponse{
		Status:        pkg.ResponseStatusOk,
		KeysProcessed: 1000,
		TotalCount:    2000,
	}, nil
}

func (s *BillingServerOkMock) GetKeyByID(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.GetKeyForOrderRequestResponse, error) {
	return &grpc.GetKeyForOrderRequestResponse{
		Status: pkg.ResponseStatusOk,
		Key:    &billing.Key{},
	}, nil
}

func (s *BillingServerOkMock) ReserveKeyForOrder(ctx context.Context, in *grpc.PlatformKeyReserveRequest, opts ...client.CallOption) (*grpc.PlatformKeyReserveResponse, error) {
	return &grpc.PlatformKeyReserveResponse{
		KeyId:  SomeMerchantId,
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) FinishRedeemKeyForOrder(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.GetKeyForOrderRequestResponse, error) {
	return &grpc.GetKeyForOrderRequestResponse{
		Status: pkg.ResponseStatusOk,
		Key:    &billing.Key{},
	}, nil
}

func (s *BillingServerOkMock) CancelRedeemKeyForOrder(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	return &grpc.EmptyResponseWithStatus{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) GetKeyProductInfo(ctx context.Context, in *grpc.GetKeyProductInfoRequest, opts ...client.CallOption) (*grpc.GetKeyProductInfoResponse, error) {
	return &grpc.GetKeyProductInfoResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkMock) ChangeCodeInOrder(ctx context.Context, in *grpc.ChangeCodeInOrderRequest, opts ...client.CallOption) (*grpc.ChangeCodeInOrderResponse, error) {
	return &grpc.ChangeCodeInOrderResponse{
		Status: pkg.ResponseStatusOk,
		Order:  &billing.Order{},
	}, nil
}

func (s *BillingServerOkMock) GetOrderPublic(
	ctx context.Context,
	in *grpc.GetOrderRequest,
	opts ...client.CallOption,
) (*grpc.GetOrderPublicResponse, error) {
	return &grpc.GetOrderPublicResponse{}, nil
}

func (s *BillingServerOkMock) GetOrderPrivate(
	ctx context.Context,
	in *grpc.GetOrderRequest,
	opts ...client.CallOption,
) (*grpc.GetOrderPrivateResponse, error) {
	return &grpc.GetOrderPrivateResponse{}, nil
}

func (s *BillingServerOkMock) FindAllOrdersPublic(
	ctx context.Context,
	in *grpc.ListOrdersRequest,
	opts ...client.CallOption,
) (*grpc.ListOrdersPublicResponse, error) {
	return &grpc.ListOrdersPublicResponse{}, nil
}

func (s *BillingServerOkMock) FindAllOrdersPrivate(
	ctx context.Context,
	in *grpc.ListOrdersRequest,
	opts ...client.CallOption,
) (*grpc.ListOrdersPrivateResponse, error) {
	return &grpc.ListOrdersPrivateResponse{}, nil
}

func (s *BillingServerOkMock) GetDashboardMainReport(
	ctx context.Context,
	in *grpc.GetDashboardMainRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardMainResponse, error) {
	return &grpc.GetDashboardMainResponse{}, nil
}
func (s *BillingServerOkMock) GetDashboardRevenueDynamicsReport(
	ctx context.Context,
	in *grpc.GetDashboardMainRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardRevenueDynamicsReportResponse, error) {
	return &grpc.GetDashboardRevenueDynamicsReportResponse{}, nil
}

func (s *BillingServerOkMock) GetDashboardBaseReport(
	ctx context.Context,
	in *grpc.GetDashboardBaseReportRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardBaseReportResponse, error) {
	return &grpc.GetDashboardBaseReportResponse{}, nil
}

func (s *BillingServerOkMock) CreatePayoutDocument(ctx context.Context, in *grpc.CreatePayoutDocumentRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) UpdatePayoutDocument(ctx context.Context, in *grpc.UpdatePayoutDocumentRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetPayoutDocuments(ctx context.Context, in *grpc.GetPayoutDocumentsRequest, opts ...client.CallOption) (*grpc.GetPayoutDocumentsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetPayoutDocumentSignUrl(ctx context.Context, in *grpc.GetPayoutDocumentSignUrlRequest, opts ...client.CallOption) (*grpc.GetPayoutDocumentSignUrlResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) UpdatePayoutDocumentSignatures(ctx context.Context, in *grpc.UpdatePayoutDocumentSignaturesRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetMerchantBalance(ctx context.Context, in *grpc.GetMerchantBalanceRequest, opts ...client.CallOption) (*grpc.GetMerchantBalanceResponse, error) {
	return &grpc.GetMerchantBalanceResponse{
		Status: pkg.ResponseStatusOk,
		Item: &billing.MerchantBalance{
			Id:             bson.NewObjectId().Hex(),
			MerchantId:     bson.NewObjectId().Hex(),
			Currency:       "RUB",
			Debit:          0,
			Credit:         0,
			RollingReserve: 0,
			Total:          0,
			CreatedAt:      ptypes.TimestampNow(),
		},
	}, nil
}

func (s *BillingServerOkMock) PayoutDocumentPdfUploaded(ctx context.Context, in *grpc.PayoutDocumentPdfUploadedRequest, opts ...client.CallOption) (*grpc.PayoutDocumentPdfUploadedResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetRoyaltyReport(ctx context.Context, in *grpc.GetRoyaltyReportRequest, opts ...client.CallOption) (*grpc.GetRoyaltyReportResponse, error) {
	return &grpc.GetRoyaltyReportResponse{
		Status:  pkg.ResponseStatusOk,
		Message: nil,
		Item:    &billing.RoyaltyReport{},
	}, nil
}

func (s *BillingServerOkMock) UnPublishKeyProduct(ctx context.Context, in *grpc.UnPublishKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) PaymentFormPlatformChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangePlatformRequest,
	opts ...client.CallOption,
) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) OrderReceipt(ctx context.Context, in *grpc.OrderReceiptRequest, opts ...client.CallOption) (*grpc.OrderReceiptResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) OrderReceiptRefund(ctx context.Context, in *grpc.OrderReceiptRequest, opts ...client.CallOption) (*grpc.OrderReceiptResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetRecommendedPriceByPriceGroup(ctx context.Context, in *grpc.RecommendedPriceRequest, opts ...client.CallOption) (*grpc.RecommendedPriceResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetRecommendedPriceByConversion(ctx context.Context, in *grpc.RecommendedPriceRequest, opts ...client.CallOption) (*grpc.RecommendedPriceResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) CheckSkuAndKeyProject(ctx context.Context, in *grpc.CheckSkuAndKeyProjectRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkMock) GetPriceGroupByRegion(ctx context.Context, in *grpc.GetPriceGroupByRegionRequest, opts ...client.CallOption) (*grpc.GetPriceGroupByRegionResponse, error) {
	return &grpc.GetPriceGroupByRegionResponse{
		Status: 200,
		Group: &billing.PriceGroup{
			Id: "some id",
		},
	}, nil
}

func (s *BillingServerOkMock) GetMerchantUsers(ctx context.Context, in *grpc.GetMerchantUsersRequest, opts ...client.CallOption) (*grpc.GetMerchantUsersResponse, error) {
	return &grpc.GetMerchantUsersResponse{
		Status: 200,
		Users: []*billing.UserRoleMerchant{
			{MerchantId: in.MerchantId, Id:SomeMerchantId},
		},
	}, nil
}
