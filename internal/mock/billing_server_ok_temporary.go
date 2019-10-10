package mock

import (
	"context"
	"github.com/globalsign/mgo/bson"
	"github.com/micro/go-micro/client"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
)

type BillingServerOkTemporaryMock struct{}

func NewBillingServerOkTemporaryMock() grpc.BillingService {
	return &BillingServerOkTemporaryMock{}
}

func (s *BillingServerOkTemporaryMock) GetProductsForOrder(
	ctx context.Context,
	in *grpc.GetProductsForOrderRequest,
	opts ...client.CallOption,
) (*grpc.ListProductsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) OrderCreateProcess(
	ctx context.Context,
	in *billing.OrderCreateRequest,
	opts ...client.CallOption,
) (*grpc.OrderCreateProcessResponse, error) {
	return &grpc.OrderCreateProcessResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Order{},
	}, nil
}

func (s *BillingServerOkTemporaryMock) PaymentFormJsonDataProcess(
	ctx context.Context,
	in *grpc.PaymentFormJsonDataRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormJsonDataResponse, error) {
	return &grpc.PaymentFormJsonDataResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) PaymentCreateProcess(
	ctx context.Context,
	in *grpc.PaymentCreateRequest,
	opts ...client.CallOption,
) (*grpc.PaymentCreateResponse, error) {
	return &grpc.PaymentCreateResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) PaymentCallbackProcess(
	ctx context.Context,
	in *grpc.PaymentNotifyRequest,
	opts ...client.CallOption,
) (*grpc.PaymentNotifyResponse, error) {
	return &grpc.PaymentNotifyResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) RebuildCache(
	ctx context.Context,
	in *grpc.EmptyRequest,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) UpdateOrder(
	ctx context.Context,
	in *billing.Order,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) UpdateMerchant(
	ctx context.Context,
	in *billing.Merchant,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) GetMerchantBy(
	ctx context.Context,
	in *grpc.GetMerchantByRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantResponse, error) {
	rsp := &grpc.GetMerchantResponse{
		Status:  pkg.ResponseStatusOk,
		Message: &grpc.ResponseErrorMessage{},
		Item:    OnboardingMerchantMock,
	}

	return rsp, nil
}

func (s *BillingServerOkTemporaryMock) ListMerchants(
	ctx context.Context,
	in *grpc.MerchantListingRequest,
	opts ...client.CallOption,
) (*grpc.MerchantListingResponse, error) {
	return &grpc.MerchantListingResponse{
		Count: 3,
		Items: []*billing.Merchant{OnboardingMerchantMock, OnboardingMerchantMock, OnboardingMerchantMock},
	}, nil
}

func (s *BillingServerOkTemporaryMock) ChangeMerchant(
	ctx context.Context,
	in *grpc.OnboardingRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantResponse, error) {
	m := &billing.Merchant{
		Company:  in.Company,
		Contacts: in.Contacts,
		Banking: &billing.MerchantBanking{
			Currency:      "RUB",
			Name:          in.Banking.Name,
			Address:       in.Banking.Address,
			AccountNumber: in.Banking.AccountNumber,
			Swift:         in.Banking.Swift,
			Details:       in.Banking.Details,
		},
		Status: pkg.MerchantStatusDraft,
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

func (s *BillingServerOkTemporaryMock) ChangeMerchantStatus(
	ctx context.Context,
	in *grpc.MerchantChangeStatusRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantStatusResponse, error) {
	return &grpc.ChangeMerchantStatusResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Merchant{Id: in.MerchantId, Status: in.Status},
	}, nil
}

func (s *BillingServerOkTemporaryMock) CreateNotification(
	ctx context.Context,
	in *grpc.NotificationRequest,
	opts ...client.CallOption,
) (*grpc.CreateNotificationResponse, error) {
	return &grpc.CreateNotificationResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Notification{},
	}, nil
}

func (s *BillingServerOkTemporaryMock) GetNotification(
	ctx context.Context,
	in *grpc.GetNotificationRequest,
	opts ...client.CallOption,
) (*billing.Notification, error) {
	return &billing.Notification{}, nil
}

func (s *BillingServerOkTemporaryMock) ListNotifications(
	ctx context.Context,
	in *grpc.ListingNotificationRequest,
	opts ...client.CallOption,
) (*grpc.Notifications, error) {
	return &grpc.Notifications{}, nil
}

func (s *BillingServerOkTemporaryMock) MarkNotificationAsRead(
	ctx context.Context,
	in *grpc.GetNotificationRequest,
	opts ...client.CallOption,
) (*billing.Notification, error) {
	return &billing.Notification{}, nil
}

func (s *BillingServerOkTemporaryMock) ListMerchantPaymentMethods(
	ctx context.Context,
	in *grpc.ListMerchantPaymentMethodsRequest,
	opts ...client.CallOption,
) (*grpc.ListingMerchantPaymentMethod, error) {
	return &grpc.ListingMerchantPaymentMethod{}, nil
}

func (s *BillingServerOkTemporaryMock) GetMerchantPaymentMethod(
	ctx context.Context,
	in *grpc.GetMerchantPaymentMethodRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantPaymentMethodResponse, error) {
	return &grpc.GetMerchantPaymentMethodResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) ChangeMerchantPaymentMethod(
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

func (s *BillingServerOkTemporaryMock) CreateRefund(
	ctx context.Context,
	in *grpc.CreateRefundRequest,
	opts ...client.CallOption,
) (*grpc.CreateRefundResponse, error) {
	return &grpc.CreateRefundResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Refund{},
	}, nil
}

func (s *BillingServerOkTemporaryMock) ListRefunds(
	ctx context.Context,
	in *grpc.ListRefundsRequest,
	opts ...client.CallOption,
) (*grpc.ListRefundsResponse, error) {
	return &grpc.ListRefundsResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) GetRefund(
	ctx context.Context,
	in *grpc.GetRefundRequest,
	opts ...client.CallOption,
) (*grpc.CreateRefundResponse, error) {
	return &grpc.CreateRefundResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Refund{},
	}, nil
}

func (s *BillingServerOkTemporaryMock) ProcessRefundCallback(
	ctx context.Context,
	in *grpc.CallbackRequest,
	opts ...client.CallOption,
) (*grpc.PaymentNotifyResponse, error) {
	return &grpc.PaymentNotifyResponse{
		Status: pkg.ResponseStatusOk,
		Error:  SomeError.Message,
	}, nil
}

func (s *BillingServerOkTemporaryMock) ChangeProject(
	ctx context.Context,
	in *billing.Project,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetProject(
	ctx context.Context,
	in *grpc.GetProjectRequest,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) DeleteProject(
	ctx context.Context,
	in *grpc.GetProjectRequest,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return &grpc.ChangeProjectResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerOkTemporaryMock) CreateToken(
	ctx context.Context,
	in *grpc.TokenRequest,
	opts ...client.CallOption,
) (*grpc.TokenResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CheckProjectRequestSignature(
	ctx context.Context,
	in *grpc.CheckProjectRequestSignatureRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CreateOrUpdateProduct(ctx context.Context, in *grpc.Product, opts ...client.CallOption) (*grpc.Product, error) {
	return Product, nil
}

func (s *BillingServerOkTemporaryMock) ListProducts(ctx context.Context, in *grpc.ListProductsRequest, opts ...client.CallOption) (*grpc.ListProductsResponse, error) {
	return &grpc.ListProductsResponse{
		Limit:  1,
		Offset: 0,
		Total:  200,
		Products: []*grpc.Product{
			Product,
		},
	}, nil
}

func (s *BillingServerOkTemporaryMock) GetProduct(ctx context.Context, in *grpc.RequestProduct, opts ...client.CallOption) (*grpc.GetProductResponse, error) {
	return GetProductResponse, nil
}

func (s *BillingServerOkTemporaryMock) DeleteProduct(ctx context.Context, in *grpc.RequestProduct, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) PaymentFormLanguageChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangeLangRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormDataChangeResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) PaymentFormPaymentAccountChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangePaymentAccountRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormDataChangeResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) ProcessBillingAddress(
	ctx context.Context,
	in *grpc.ProcessBillingAddressRequest,
	opts ...client.CallOption,
) (*grpc.ProcessBillingAddressResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) ChangeMerchantData(
	ctx context.Context,
	in *grpc.ChangeMerchantDataRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantDataResponse, error) {
	return &grpc.ChangeMerchantDataResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) SetMerchantS3Agreement(
	ctx context.Context,
	in *grpc.SetMerchantS3AgreementRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantDataResponse, error) {
	return &grpc.ChangeMerchantDataResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) FindAllOrders(ctx context.Context, in *grpc.ListOrdersRequest, opts ...client.CallOption) (*billing.OrderPaginate, error) {
	return &billing.OrderPaginate{}, nil
}

func (s *BillingServerOkTemporaryMock) ListProjects(ctx context.Context, in *grpc.ListProjectsRequest, opts ...client.CallOption) (*grpc.ListProjectsResponse, error) {
	return &grpc.ListProjectsResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) GetOrder(ctx context.Context, in *grpc.GetOrderRequest, opts ...client.CallOption) (*billing.Order, error) {
	return &billing.Order{}, nil
}

func (s *BillingServerOkTemporaryMock) IsOrderCanBePaying(
	ctx context.Context,
	in *grpc.IsOrderCanBePayingRequest,
	opts ...client.CallOption,
) (*grpc.IsOrderCanBePayingResponse, error) {
	return &grpc.IsOrderCanBePayingResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Order{},
	}, nil
}

func (s *BillingServerOkTemporaryMock) GetCountry(ctx context.Context, in *billing.GetCountryRequest, opts ...client.CallOption) (*billing.Country, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UpdateCountry(ctx context.Context, in *billing.Country, opts ...client.CallOption) (*billing.Country, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetPriceGroup(ctx context.Context, in *billing.GetPriceGroupRequest, opts ...client.CallOption) (*billing.PriceGroup, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UpdatePriceGroup(ctx context.Context, in *billing.PriceGroup, opts ...client.CallOption) (*billing.PriceGroup, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) SetUserNotifySales(ctx context.Context, in *grpc.SetUserNotifyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) SetUserNotifyNewRegion(ctx context.Context, in *grpc.SetUserNotifyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}
func (s *BillingServerOkTemporaryMock) GetCountriesList(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*billing.CountriesList, error) {
	panic("implement me")
}
func (s *BillingServerOkTemporaryMock) GetPaymentChannelCostSystem(ctx context.Context, in *billing.PaymentChannelCostSystemRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) SetPaymentChannelCostSystem(ctx context.Context, in *billing.PaymentChannelCostSystem, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeletePaymentChannelCostSystem(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchantRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) SetPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchant, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeletePaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetMoneyBackCostSystem(ctx context.Context, in *billing.MoneyBackCostSystemRequest, opts ...client.CallOption) (*grpc.MoneyBackCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) SetMoneyBackCostSystem(ctx context.Context, in *billing.MoneyBackCostSystem, opts ...client.CallOption) (*grpc.MoneyBackCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeleteMoneyBackCostSystem(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchantRequest, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) SetMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchant, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeleteMoneyBackCostMerchant(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetAllPaymentChannelCostSystem(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemListResponse, error) {
	return &grpc.PaymentChannelCostSystemListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkTemporaryMock) GetAllPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchantListRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantListResponse, error) {
	return &grpc.PaymentChannelCostMerchantListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkTemporaryMock) GetAllMoneyBackCostSystem(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.MoneyBackCostSystemListResponse, error) {
	return &grpc.MoneyBackCostSystemListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkTemporaryMock) GetAllMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchantListRequest, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantListResponse, error) {
	return &grpc.MoneyBackCostMerchantListResponse{
		Status: pkg.ResponseStatusOk,
	}, nil
}

func (s *BillingServerOkTemporaryMock) CreateOrUpdatePaymentMethodTestSettings(ctx context.Context, in *grpc.ChangePaymentMethodParamsRequest, opts ...client.CallOption) (*grpc.ChangePaymentMethodParamsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeletePaymentMethodTestSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.ChangePaymentMethodParamsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) FindByZipCode(
	ctx context.Context,
	in *grpc.FindByZipCodeRequest,
	opts ...client.CallOption,
) (*grpc.FindByZipCodeResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CreateOrUpdatePaymentMethod(
	ctx context.Context,
	in *billing.PaymentMethod,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CreateOrUpdatePaymentMethodProductionSettings(
	ctx context.Context,
	in *grpc.ChangePaymentMethodParamsRequest,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodParamsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) DeletePaymentMethodProductionSettings(
	ctx context.Context,
	in *grpc.GetPaymentMethodSettingsRequest,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodParamsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CreateAccountingEntry(ctx context.Context, in *grpc.CreateAccountingEntryRequest, opts ...client.CallOption) (*grpc.CreateAccountingEntryResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) CreateRoyaltyReport(ctx context.Context, in *grpc.CreateRoyaltyReportRequest, opts ...client.CallOption) (*grpc.CreateRoyaltyReportRequest, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) ListRoyaltyReports(ctx context.Context, in *grpc.ListRoyaltyReportsRequest, opts ...client.CallOption) (*grpc.ListRoyaltyReportsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) ChangeRoyaltyReportStatus(ctx context.Context, in *grpc.CreateRoyaltyReportRequest, opts ...client.CallOption) (*grpc.CreateRoyaltyReportRequest, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) ListRoyaltyReportOrders(ctx context.Context, in *grpc.ListRoyaltyReportOrdersRequest, opts ...client.CallOption) (*grpc.TransactionsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetVatReportsDashboard(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.VatReportsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetVatReportsForCountry(ctx context.Context, in *grpc.VatReportsRequest, opts ...client.CallOption) (*grpc.VatReportsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetVatReportTransactions(ctx context.Context, in *grpc.VatTransactionsRequest, opts ...client.CallOption) (*grpc.TransactionsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) CalcAnnualTurnovers(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}
func (s *BillingServerOkTemporaryMock) ProcessVatReports(ctx context.Context, in *grpc.ProcessVatReportsRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UpdateVatReportStatus(ctx context.Context, in *grpc.UpdateVatReportStatusRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UpdateProductPrices(
	ctx context.Context,
	in *grpc.UpdateProductPricesRequest,
	opts ...client.CallOption,
) (*grpc.ResponseError, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetProductPrices(
	ctx context.Context,
	in *grpc.RequestProduct,
	opts ...client.CallOption,
) (*grpc.ProductPricesResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetPriceGroupRecommendedPrice(
	ctx context.Context,
	in *grpc.RecommendedPriceRequest,
	opts ...client.CallOption,
) (*grpc.RecommendedPriceResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetPriceGroupCurrencyByRegion(
	ctx context.Context,
	in *grpc.PriceGroupByRegionRequest,
	opts ...client.CallOption,
) (*grpc.PriceGroupCurrenciesResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetPriceGroupCurrencies(
	ctx context.Context,
	in *grpc.EmptyRequest,
	opts ...client.CallOption,
) (*grpc.PriceGroupCurrenciesResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetPriceGroupByCountry(
	ctx context.Context,
	in *grpc.PriceGroupByCountryRequest,
	opts ...client.CallOption,
) (*billing.PriceGroup, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetPaymentMethodProductionSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.GetPaymentMethodSettingsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetPaymentMethodTestSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.GetPaymentMethodSettingsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) ChangeRoyaltyReport(ctx context.Context, in *grpc.ChangeRoyaltyReportRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) AutoAcceptRoyaltyReports(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetUserProfile(
	ctx context.Context,
	in *grpc.GetUserProfileRequest,
	opts ...client.CallOption,
) (*grpc.GetUserProfileResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CreateOrUpdateUserProfile(
	ctx context.Context,
	in *grpc.UserProfile,
	opts ...client.CallOption,
) (*grpc.GetUserProfileResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) ConfirmUserEmail(
	ctx context.Context,
	in *grpc.ConfirmUserEmailRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) CreatePageReview(
	ctx context.Context,
	in *grpc.CreatePageReviewRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) MerchantReviewRoyaltyReport(ctx context.Context, in *grpc.MerchantReviewRoyaltyReportRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetMerchantAgreementSignUrl(
	ctx context.Context,
	in *grpc.GetMerchantAgreementSignUrlRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantAgreementSignUrlResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetMerchantOnboardingCompleteData(
	ctx context.Context,
	in *grpc.SetMerchantS3AgreementRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantOnboardingCompleteDataResponse, error) {
	return nil, SomeError
}

func (s *BillingServerOkTemporaryMock) GetMerchantTariffRates(
	ctx context.Context,
	in *grpc.GetMerchantTariffRatesRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantTariffRatesResponse, error) {
	return &grpc.GetMerchantTariffRatesResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) SetMerchantTariffRates(
	ctx context.Context,
	in *grpc.SetMerchantTariffRatesRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) CreateOrUpdateKeyProduct(ctx context.Context, in *grpc.CreateOrUpdateKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetKeyProducts(ctx context.Context, in *grpc.ListKeyProductsRequest, opts ...client.CallOption) (*grpc.ListKeyProductsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetKeyProduct(ctx context.Context, in *grpc.RequestKeyProductMerchant, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeleteKeyProduct(ctx context.Context, in *grpc.RequestKeyProductMerchant, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) PublishKeyProduct(ctx context.Context, in *grpc.PublishKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetKeyProductsForOrder(ctx context.Context, in *grpc.GetKeyProductsForOrderRequest, opts ...client.CallOption) (*grpc.ListKeyProductsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetPlatforms(ctx context.Context, in *grpc.ListPlatformsRequest, opts ...client.CallOption) (*grpc.ListPlatformsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) DeletePlatformFromProduct(ctx context.Context, in *grpc.RemovePlatformRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetAvailableKeysCount(ctx context.Context, in *grpc.GetPlatformKeyCountRequest, opts ...client.CallOption) (*grpc.GetPlatformKeyCountResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UploadKeysFile(ctx context.Context, in *grpc.PlatformKeysFileRequest, opts ...client.CallOption) (*grpc.PlatformKeysFileResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetKeyByID(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.GetKeyForOrderRequestResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) ReserveKeyForOrder(ctx context.Context, in *grpc.PlatformKeyReserveRequest, opts ...client.CallOption) (*grpc.PlatformKeyReserveResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) FinishRedeemKeyForOrder(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.GetKeyForOrderRequestResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) CancelRedeemKeyForOrder(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetKeyProductInfo(ctx context.Context, in *grpc.GetKeyProductInfoRequest, opts ...client.CallOption) (*grpc.GetKeyProductInfoResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) ChangeCodeInOrder(ctx context.Context, in *grpc.ChangeCodeInOrderRequest, opts ...client.CallOption) (*grpc.ChangeCodeInOrderResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetDashboardMainReport(
	ctx context.Context,
	in *grpc.GetDashboardMainRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardMainResponse, error) {
	return &grpc.GetDashboardMainResponse{}, nil
}
func (s *BillingServerOkTemporaryMock) GetDashboardRevenueDynamicsReport(
	ctx context.Context,
	in *grpc.GetDashboardMainRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardRevenueDynamicsReportResponse, error) {
	return &grpc.GetDashboardRevenueDynamicsReportResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) GetDashboardBaseReport(
	ctx context.Context,
	in *grpc.GetDashboardBaseReportRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardBaseReportResponse, error) {
	return &grpc.GetDashboardBaseReportResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) GetOrderPublic(
	ctx context.Context,
	in *grpc.GetOrderRequest,
	opts ...client.CallOption,
) (*grpc.GetOrderPublicResponse, error) {
	return &grpc.GetOrderPublicResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) GetOrderPrivate(
	ctx context.Context,
	in *grpc.GetOrderRequest,
	opts ...client.CallOption,
) (*grpc.GetOrderPrivateResponse, error) {
	return &grpc.GetOrderPrivateResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) FindAllOrdersPublic(
	ctx context.Context,
	in *grpc.ListOrdersRequest,
	opts ...client.CallOption,
) (*grpc.ListOrdersPublicResponse, error) {
	return &grpc.ListOrdersPublicResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) FindAllOrdersPrivate(
	ctx context.Context,
	in *grpc.ListOrdersRequest,
	opts ...client.CallOption,
) (*grpc.ListOrdersPrivateResponse, error) {
	return &grpc.ListOrdersPrivateResponse{}, nil
}

func (s *BillingServerOkTemporaryMock) CreatePayoutDocument(ctx context.Context, in *grpc.CreatePayoutDocumentRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UpdatePayoutDocument(ctx context.Context, in *grpc.UpdatePayoutDocumentRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetPayoutDocuments(ctx context.Context, in *grpc.GetPayoutDocumentsRequest, opts ...client.CallOption) (*grpc.GetPayoutDocumentsResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetPayoutDocumentSignUrl(ctx context.Context, in *grpc.GetPayoutDocumentSignUrlRequest, opts ...client.CallOption) (*grpc.GetPayoutDocumentSignUrlResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UpdatePayoutDocumentSignatures(ctx context.Context, in *grpc.UpdatePayoutDocumentSignaturesRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetMerchantBalance(ctx context.Context, in *grpc.GetMerchantBalanceRequest, opts ...client.CallOption) (*grpc.GetMerchantBalanceResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) PayoutDocumentPdfUploaded(ctx context.Context, in *grpc.PayoutDocumentPdfUploadedRequest, opts ...client.CallOption) (*grpc.PayoutDocumentPdfUploadedResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetRoyaltyReport(ctx context.Context, in *grpc.GetRoyaltyReportRequest, opts ...client.CallOption) (*grpc.GetRoyaltyReportResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) UnPublishKeyProduct(ctx context.Context, in *grpc.UnPublishKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) PaymentFormPlatformChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangePlatformRequest,
	opts ...client.CallOption,
) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) OrderReceipt(ctx context.Context, in *grpc.OrderReceiptRequest, opts ...client.CallOption) (*grpc.OrderReceiptResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) OrderReceiptRefund(ctx context.Context, in *grpc.OrderReceiptRequest, opts ...client.CallOption) (*grpc.OrderReceiptResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetRecommendedPriceByPriceGroup(ctx context.Context, in *grpc.RecommendedPriceRequest, opts ...client.CallOption) (*grpc.RecommendedPriceResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetRecommendedPriceByConversion(ctx context.Context, in *grpc.RecommendedPriceRequest, opts ...client.CallOption) (*grpc.RecommendedPriceResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) CheckSkuAndKeyProject(ctx context.Context, in *grpc.CheckSkuAndKeyProjectRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetPriceGroupByRegion(ctx context.Context, in *grpc.GetPriceGroupByRegionRequest, opts ...client.CallOption) (*grpc.GetPriceGroupByRegionResponse, error) {
	panic("implement me")
}

func (s *BillingServerOkTemporaryMock) GetMerchantUsers(ctx context.Context, in *grpc.GetMerchantUsersRequest, opts ...client.CallOption) (*grpc.GetMerchantUsersResponse, error) {
	panic("implement me")
}
