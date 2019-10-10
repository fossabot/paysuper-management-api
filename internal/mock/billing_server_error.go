package mock

import (
	"context"
	"errors"
	"github.com/globalsign/mgo/bson"
	"github.com/micro/go-micro/client"
	"github.com/paysuper/paysuper-billing-server/pkg"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/billing"
	"github.com/paysuper/paysuper-billing-server/pkg/proto/grpc"
	"net/http"
)

type BillingServerErrorMock struct{}

func NewBillingServerErrorMock() grpc.BillingService {
	return &BillingServerErrorMock{}
}

func (s *BillingServerErrorMock) GetProductsForOrder(
	ctx context.Context,
	in *grpc.GetProductsForOrderRequest,
	opts ...client.CallOption,
) (*grpc.ListProductsResponse, error) {
	return &grpc.ListProductsResponse{}, nil
}

func (s *BillingServerErrorMock) OrderCreateProcess(
	ctx context.Context,
	in *billing.OrderCreateRequest,
	opts ...client.CallOption,
) (*grpc.OrderCreateProcessResponse, error) {
	return &grpc.OrderCreateProcessResponse{
		Status: pkg.ResponseStatusOk,
		Item:   &billing.Order{},
	}, nil
}

func (s *BillingServerErrorMock) PaymentFormJsonDataProcess(
	ctx context.Context,
	in *grpc.PaymentFormJsonDataRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormJsonDataResponse, error) {
	return &grpc.PaymentFormJsonDataResponse{}, nil
}

func (s *BillingServerErrorMock) PaymentCreateProcess(
	ctx context.Context,
	in *grpc.PaymentCreateRequest,
	opts ...client.CallOption,
) (*grpc.PaymentCreateResponse, error) {
	return &grpc.PaymentCreateResponse{}, nil
}

func (s *BillingServerErrorMock) PaymentCallbackProcess(
	ctx context.Context,
	in *grpc.PaymentNotifyRequest,
	opts ...client.CallOption,
) (*grpc.PaymentNotifyResponse, error) {
	return &grpc.PaymentNotifyResponse{}, nil
}

func (s *BillingServerErrorMock) RebuildCache(
	ctx context.Context,
	in *grpc.EmptyRequest,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerErrorMock) UpdateOrder(
	ctx context.Context,
	in *billing.Order,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerErrorMock) UpdateMerchant(
	ctx context.Context,
	in *billing.Merchant,
	opts ...client.CallOption,
) (*grpc.EmptyResponse, error) {
	return &grpc.EmptyResponse{}, nil
}

func (s *BillingServerErrorMock) GetMerchantBy(
	ctx context.Context,
	in *grpc.GetMerchantByRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantResponse, error) {
	return &grpc.GetMerchantResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ListMerchants(
	ctx context.Context,
	in *grpc.MerchantListingRequest,
	opts ...client.CallOption,
) (*grpc.MerchantListingResponse, error) {
	return &grpc.MerchantListingResponse{}, nil
}

func (s *BillingServerErrorMock) ChangeMerchant(
	ctx context.Context,
	in *grpc.OnboardingRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantResponse, error) {
	return &grpc.ChangeMerchantResponse{
		Status:  http.StatusBadRequest,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ChangeMerchantStatus(
	ctx context.Context,
	in *grpc.MerchantChangeStatusRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantStatusResponse, error) {
	return &grpc.ChangeMerchantStatusResponse{
		Status:  http.StatusBadRequest,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CreateNotification(
	ctx context.Context,
	in *grpc.NotificationRequest,
	opts ...client.CallOption,
) (*grpc.CreateNotificationResponse, error) {
	return &grpc.CreateNotificationResponse{
		Status:  http.StatusBadRequest,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetNotification(
	ctx context.Context,
	in *grpc.GetNotificationRequest,
	opts ...client.CallOption,
) (*billing.Notification, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) ListNotifications(
	ctx context.Context,
	in *grpc.ListingNotificationRequest,
	opts ...client.CallOption,
) (*grpc.Notifications, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) MarkNotificationAsRead(
	ctx context.Context,
	in *grpc.GetNotificationRequest,
	opts ...client.CallOption,
) (*billing.Notification, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) ListMerchantPaymentMethods(
	ctx context.Context,
	in *grpc.ListMerchantPaymentMethodsRequest,
	opts ...client.CallOption,
) (*grpc.ListingMerchantPaymentMethod, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetMerchantPaymentMethod(
	ctx context.Context,
	in *grpc.GetMerchantPaymentMethodRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantPaymentMethodResponse, error) {
	return &grpc.GetMerchantPaymentMethodResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ChangeMerchantPaymentMethod(
	ctx context.Context,
	in *grpc.MerchantPaymentMethodRequest,
	opts ...client.CallOption,
) (*grpc.MerchantPaymentMethodResponse, error) {
	return &grpc.MerchantPaymentMethodResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CreateRefund(
	ctx context.Context,
	in *grpc.CreateRefundRequest,
	opts ...client.CallOption,
) (*grpc.CreateRefundResponse, error) {
	return &grpc.CreateRefundResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ListRefunds(
	ctx context.Context,
	in *grpc.ListRefundsRequest,
	opts ...client.CallOption,
) (*grpc.ListRefundsResponse, error) {
	return &grpc.ListRefundsResponse{}, nil
}

func (s *BillingServerErrorMock) GetRefund(
	ctx context.Context,
	in *grpc.GetRefundRequest,
	opts ...client.CallOption,
) (*grpc.CreateRefundResponse, error) {
	return &grpc.CreateRefundResponse{
		Status:  pkg.ResponseStatusNotFound,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ProcessRefundCallback(
	ctx context.Context,
	in *grpc.CallbackRequest,
	opts ...client.CallOption,
) (*grpc.PaymentNotifyResponse, error) {
	return &grpc.PaymentNotifyResponse{
		Status: pkg.ResponseStatusNotFound,
		Error:  SomeError.Message,
	}, nil
}

func (s *BillingServerErrorMock) PaymentFormLanguageChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangeLangRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormDataChangeResponse, error) {
	return &grpc.PaymentFormDataChangeResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) PaymentFormPaymentAccountChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangePaymentAccountRequest,
	opts ...client.CallOption,
) (*grpc.PaymentFormDataChangeResponse, error) {
	return &grpc.PaymentFormDataChangeResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ProcessBillingAddress(
	ctx context.Context,
	in *grpc.ProcessBillingAddressRequest,
	opts ...client.CallOption,
) (*grpc.ProcessBillingAddressResponse, error) {
	return &grpc.ProcessBillingAddressResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ChangeMerchantData(
	ctx context.Context,
	in *grpc.ChangeMerchantDataRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantDataResponse, error) {
	return &grpc.ChangeMerchantDataResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) SetMerchantS3Agreement(
	ctx context.Context,
	in *grpc.SetMerchantS3AgreementRequest,
	opts ...client.CallOption,
) (*grpc.ChangeMerchantDataResponse, error) {
	return &grpc.ChangeMerchantDataResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ChangeProject(
	ctx context.Context,
	in *billing.Project,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return &grpc.ChangeProjectResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetProject(
	ctx context.Context,
	in *grpc.GetProjectRequest,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	if in.ProjectId == SomeMerchantId {
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

	return &grpc.ChangeProjectResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) DeleteProject(
	ctx context.Context,
	in *grpc.GetProjectRequest,
	opts ...client.CallOption,
) (*grpc.ChangeProjectResponse, error) {
	return &grpc.ChangeProjectResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CreateToken(
	ctx context.Context,
	in *grpc.TokenRequest,
	opts ...client.CallOption,
) (*grpc.TokenResponse, error) {
	return &grpc.TokenResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CheckProjectRequestSignature(
	ctx context.Context,
	in *grpc.CheckProjectRequestSignatureRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CreateOrUpdateProduct(ctx context.Context, in *grpc.Product, opts ...client.CallOption) (*grpc.Product, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) ListProducts(ctx context.Context, in *grpc.ListProductsRequest, opts ...client.CallOption) (*grpc.ListProductsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetProduct(ctx context.Context, in *grpc.RequestProduct, opts ...client.CallOption) (*grpc.GetProductResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) DeleteProduct(ctx context.Context, in *grpc.RequestProduct, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) FindAllOrders(ctx context.Context, in *grpc.ListOrdersRequest, opts ...client.CallOption) (*billing.OrderPaginate, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) ListProjects(ctx context.Context, in *grpc.ListProjectsRequest, opts ...client.CallOption) (*grpc.ListProjectsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetOrder(ctx context.Context, in *grpc.GetOrderRequest, opts ...client.CallOption) (*billing.Order, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) IsOrderCanBePaying(
	ctx context.Context,
	in *grpc.IsOrderCanBePayingRequest,
	opts ...client.CallOption,
) (*grpc.IsOrderCanBePayingResponse, error) {
	return &grpc.IsOrderCanBePayingResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetCountry(ctx context.Context, in *billing.GetCountryRequest, opts ...client.CallOption) (*billing.Country, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) UpdateCountry(ctx context.Context, in *billing.Country, opts ...client.CallOption) (*billing.Country, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPriceGroup(ctx context.Context, in *billing.GetPriceGroupRequest, opts ...client.CallOption) (*billing.PriceGroup, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) UpdatePriceGroup(ctx context.Context, in *billing.PriceGroup, opts ...client.CallOption) (*billing.PriceGroup, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) SetUserNotifySales(ctx context.Context, in *grpc.SetUserNotifyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) SetUserNotifyNewRegion(ctx context.Context, in *grpc.SetUserNotifyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetCountriesList(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*billing.CountriesList, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPaymentChannelCostSystem(ctx context.Context, in *billing.PaymentChannelCostSystemRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) SetPaymentChannelCostSystem(ctx context.Context, in *billing.PaymentChannelCostSystem, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) DeletePaymentChannelCostSystem(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchantRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) SetPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchant, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) DeletePaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetMoneyBackCostSystem(ctx context.Context, in *billing.MoneyBackCostSystemRequest, opts ...client.CallOption) (*grpc.MoneyBackCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) SetMoneyBackCostSystem(ctx context.Context, in *billing.MoneyBackCostSystem, opts ...client.CallOption) (*grpc.MoneyBackCostSystemResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) DeleteMoneyBackCostSystem(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchantRequest, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) SetMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchant, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) DeleteMoneyBackCostMerchant(ctx context.Context, in *billing.PaymentCostDeleteRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetAllPaymentChannelCostSystem(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostSystemListResponse, error) {
	return nil, errors.New("Some error")
}

func (s *BillingServerErrorMock) GetAllPaymentChannelCostMerchant(ctx context.Context, in *billing.PaymentChannelCostMerchantListRequest, opts ...client.CallOption) (*grpc.PaymentChannelCostMerchantListResponse, error) {
	return nil, errors.New("Some error")
}

func (s *BillingServerErrorMock) GetAllMoneyBackCostSystem(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.MoneyBackCostSystemListResponse, error) {
	return nil, errors.New("Some error")
}

func (s *BillingServerErrorMock) GetAllMoneyBackCostMerchant(ctx context.Context, in *billing.MoneyBackCostMerchantListRequest, opts ...client.CallOption) (*grpc.MoneyBackCostMerchantListResponse, error) {
	return nil, errors.New("Some error")
}

func (s *BillingServerErrorMock) CreateOrUpdatePaymentMethodTestSettings(ctx context.Context, in *grpc.ChangePaymentMethodParamsRequest, opts ...client.CallOption) (*grpc.ChangePaymentMethodParamsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) DeletePaymentMethodTestSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.ChangePaymentMethodParamsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) FindByZipCode(
	ctx context.Context,
	in *grpc.FindByZipCodeRequest,
	opts ...client.CallOption,
) (*grpc.FindByZipCodeResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) CreateOrUpdatePaymentMethod(
	ctx context.Context,
	in *billing.PaymentMethod,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) CreateOrUpdatePaymentMethodProductionSettings(
	ctx context.Context,
	in *grpc.ChangePaymentMethodParamsRequest,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodParamsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) DeletePaymentMethodProductionSettings(
	ctx context.Context,
	in *grpc.GetPaymentMethodSettingsRequest,
	opts ...client.CallOption,
) (*grpc.ChangePaymentMethodParamsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) CreateAccountingEntry(ctx context.Context, in *grpc.CreateAccountingEntryRequest, opts ...client.CallOption) (*grpc.CreateAccountingEntryResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) CreateRoyaltyReport(ctx context.Context, in *grpc.CreateRoyaltyReportRequest, opts ...client.CallOption) (*grpc.CreateRoyaltyReportRequest, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) ListRoyaltyReports(ctx context.Context, in *grpc.ListRoyaltyReportsRequest, opts ...client.CallOption) (*grpc.ListRoyaltyReportsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) ChangeRoyaltyReportStatus(ctx context.Context, in *grpc.CreateRoyaltyReportRequest, opts ...client.CallOption) (*grpc.CreateRoyaltyReportRequest, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) ListRoyaltyReportOrders(ctx context.Context, in *grpc.ListRoyaltyReportOrdersRequest, opts ...client.CallOption) (*grpc.TransactionsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetVatReportsDashboard(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.VatReportsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetVatReportsForCountry(ctx context.Context, in *grpc.VatReportsRequest, opts ...client.CallOption) (*grpc.VatReportsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetVatReportTransactions(ctx context.Context, in *grpc.VatTransactionsRequest, opts ...client.CallOption) (*grpc.TransactionsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) CalcAnnualTurnovers(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) ProcessVatReports(ctx context.Context, in *grpc.ProcessVatReportsRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) UpdateVatReportStatus(ctx context.Context, in *grpc.UpdateVatReportStatusRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPriceGroupByCountry(
	ctx context.Context,
	in *grpc.PriceGroupByCountryRequest,
	opts ...client.CallOption,
) (*billing.PriceGroup, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetPriceGroupCurrencies(
	ctx context.Context,
	in *grpc.EmptyRequest,
	opts ...client.CallOption,
) (*grpc.PriceGroupCurrenciesResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetPriceGroupCurrencyByRegion(
	ctx context.Context,
	in *grpc.PriceGroupByRegionRequest,
	opts ...client.CallOption,
) (*grpc.PriceGroupCurrenciesResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetPriceGroupRecommendedPrice(
	ctx context.Context,
	in *grpc.RecommendedPriceRequest,
	opts ...client.CallOption,
) (*grpc.RecommendedPriceResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetProductPrices(
	ctx context.Context,
	in *grpc.RequestProduct,
	opts ...client.CallOption,
) (*grpc.ProductPricesResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) UpdateProductPrices(
	ctx context.Context,
	in *grpc.UpdateProductPricesRequest,
	opts ...client.CallOption,
) (*grpc.ResponseError, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetPaymentMethodProductionSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.GetPaymentMethodSettingsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) GetPaymentMethodTestSettings(ctx context.Context, in *grpc.GetPaymentMethodSettingsRequest, opts ...client.CallOption) (*grpc.GetPaymentMethodSettingsResponse, error) {
	return nil, SomeError
}

func (s *BillingServerErrorMock) ChangeRoyaltyReport(ctx context.Context, in *grpc.ChangeRoyaltyReportRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) AutoAcceptRoyaltyReports(ctx context.Context, in *grpc.EmptyRequest, opts ...client.CallOption) (*grpc.EmptyResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetUserProfile(
	ctx context.Context,
	in *grpc.GetUserProfileRequest,
	opts ...client.CallOption,
) (*grpc.GetUserProfileResponse, error) {
	return &grpc.GetUserProfileResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CreateOrUpdateUserProfile(
	ctx context.Context,
	in *grpc.UserProfile,
	opts ...client.CallOption,
) (*grpc.GetUserProfileResponse, error) {
	return &grpc.GetUserProfileResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ConfirmUserEmail(
	ctx context.Context,
	in *grpc.ConfirmUserEmailRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CreatePageReview(
	ctx context.Context,
	in *grpc.CreatePageReviewRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) MerchantReviewRoyaltyReport(ctx context.Context, in *grpc.MerchantReviewRoyaltyReportRequest, opts ...client.CallOption) (*grpc.ResponseError, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetMerchantAgreementSignUrl(
	ctx context.Context,
	in *grpc.GetMerchantAgreementSignUrlRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantAgreementSignUrlResponse, error) {
	return &grpc.GetMerchantAgreementSignUrlResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetMerchantOnboardingCompleteData(
	ctx context.Context,
	in *grpc.SetMerchantS3AgreementRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantOnboardingCompleteDataResponse, error) {
	return &grpc.GetMerchantOnboardingCompleteDataResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetMerchantTariffRates(
	ctx context.Context,
	in *grpc.GetMerchantTariffRatesRequest,
	opts ...client.CallOption,
) (*grpc.GetMerchantTariffRatesResponse, error) {
	return &grpc.GetMerchantTariffRatesResponse{}, nil
}

func (s *BillingServerErrorMock) SetMerchantTariffRates(
	ctx context.Context,
	in *grpc.SetMerchantTariffRatesRequest,
	opts ...client.CallOption,
) (*grpc.CheckProjectRequestSignatureResponse, error) {
	return &grpc.CheckProjectRequestSignatureResponse{}, nil
}

func (s *BillingServerErrorMock) CreateOrUpdateKeyProduct(ctx context.Context, in *grpc.CreateOrUpdateKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	return &grpc.KeyProductResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetKeyProducts(ctx context.Context, in *grpc.ListKeyProductsRequest, opts ...client.CallOption) (*grpc.ListKeyProductsResponse, error) {
	return &grpc.ListKeyProductsResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetKeyProduct(ctx context.Context, in *grpc.RequestKeyProductMerchant, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	return &grpc.KeyProductResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) DeleteKeyProduct(ctx context.Context, in *grpc.RequestKeyProductMerchant, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	return &grpc.EmptyResponseWithStatus{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) PublishKeyProduct(ctx context.Context, in *grpc.PublishKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	return &grpc.KeyProductResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetKeyProductsForOrder(ctx context.Context, in *grpc.GetKeyProductsForOrderRequest, opts ...client.CallOption) (*grpc.ListKeyProductsResponse, error) {
	return &grpc.ListKeyProductsResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetPlatforms(ctx context.Context, in *grpc.ListPlatformsRequest, opts ...client.CallOption) (*grpc.ListPlatformsResponse, error) {
	return &grpc.ListPlatformsResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) DeletePlatformFromProduct(ctx context.Context, in *grpc.RemovePlatformRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	return &grpc.EmptyResponseWithStatus{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetAvailableKeysCount(ctx context.Context, in *grpc.GetPlatformKeyCountRequest, opts ...client.CallOption) (*grpc.GetPlatformKeyCountResponse, error) {
	return &grpc.GetPlatformKeyCountResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) UploadKeysFile(ctx context.Context, in *grpc.PlatformKeysFileRequest, opts ...client.CallOption) (*grpc.PlatformKeysFileResponse, error) {
	return &grpc.PlatformKeysFileResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetKeyByID(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.GetKeyForOrderRequestResponse, error) {
	return &grpc.GetKeyForOrderRequestResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ReserveKeyForOrder(ctx context.Context, in *grpc.PlatformKeyReserveRequest, opts ...client.CallOption) (*grpc.PlatformKeyReserveResponse, error) {
	return &grpc.PlatformKeyReserveResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) FinishRedeemKeyForOrder(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.GetKeyForOrderRequestResponse, error) {
	return &grpc.GetKeyForOrderRequestResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) CancelRedeemKeyForOrder(ctx context.Context, in *grpc.KeyForOrderRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	return &grpc.EmptyResponseWithStatus{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) GetKeyProductInfo(ctx context.Context, in *grpc.GetKeyProductInfoRequest, opts ...client.CallOption) (*grpc.GetKeyProductInfoResponse, error) {
	return &grpc.GetKeyProductInfoResponse{
		Status:  pkg.ResponseStatusBadData,
		Message: SomeError,
	}, nil
}

func (s *BillingServerErrorMock) ChangeCodeInOrder(ctx context.Context, in *grpc.ChangeCodeInOrderRequest, opts ...client.CallOption) (*grpc.ChangeCodeInOrderResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetDashboardMainReport(
	ctx context.Context,
	in *grpc.GetDashboardMainRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardMainResponse, error) {
	return &grpc.GetDashboardMainResponse{}, nil
}
func (s *BillingServerErrorMock) GetDashboardRevenueDynamicsReport(
	ctx context.Context,
	in *grpc.GetDashboardMainRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardRevenueDynamicsReportResponse, error) {
	return &grpc.GetDashboardRevenueDynamicsReportResponse{}, nil
}

func (s *BillingServerErrorMock) GetDashboardBaseReport(
	ctx context.Context,
	in *grpc.GetDashboardBaseReportRequest,
	opts ...client.CallOption,
) (*grpc.GetDashboardBaseReportResponse, error) {
	return &grpc.GetDashboardBaseReportResponse{}, nil
}

func (s *BillingServerErrorMock) GetOrderPublic(
	ctx context.Context,
	in *grpc.GetOrderRequest,
	opts ...client.CallOption,
) (*grpc.GetOrderPublicResponse, error) {
	return &grpc.GetOrderPublicResponse{}, nil
}

func (s *BillingServerErrorMock) GetOrderPrivate(
	ctx context.Context,
	in *grpc.GetOrderRequest,
	opts ...client.CallOption,
) (*grpc.GetOrderPrivateResponse, error) {
	return &grpc.GetOrderPrivateResponse{}, nil
}

func (s *BillingServerErrorMock) FindAllOrdersPublic(
	ctx context.Context,
	in *grpc.ListOrdersRequest,
	opts ...client.CallOption,
) (*grpc.ListOrdersPublicResponse, error) {
	return &grpc.ListOrdersPublicResponse{}, nil
}

func (s *BillingServerErrorMock) FindAllOrdersPrivate(
	ctx context.Context,
	in *grpc.ListOrdersRequest,
	opts ...client.CallOption,
) (*grpc.ListOrdersPrivateResponse, error) {
	return &grpc.ListOrdersPrivateResponse{}, nil
}

func (s *BillingServerErrorMock) CreatePayoutDocument(ctx context.Context, in *grpc.CreatePayoutDocumentRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) UpdatePayoutDocument(ctx context.Context, in *grpc.UpdatePayoutDocumentRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPayoutDocuments(ctx context.Context, in *grpc.GetPayoutDocumentsRequest, opts ...client.CallOption) (*grpc.GetPayoutDocumentsResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPayoutDocumentSignUrl(ctx context.Context, in *grpc.GetPayoutDocumentSignUrlRequest, opts ...client.CallOption) (*grpc.GetPayoutDocumentSignUrlResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) UpdatePayoutDocumentSignatures(ctx context.Context, in *grpc.UpdatePayoutDocumentSignaturesRequest, opts ...client.CallOption) (*grpc.PayoutDocumentResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetMerchantBalance(ctx context.Context, in *grpc.GetMerchantBalanceRequest, opts ...client.CallOption) (*grpc.GetMerchantBalanceResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) PayoutDocumentPdfUploaded(ctx context.Context, in *grpc.PayoutDocumentPdfUploadedRequest, opts ...client.CallOption) (*grpc.PayoutDocumentPdfUploadedResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetRoyaltyReport(ctx context.Context, in *grpc.GetRoyaltyReportRequest, opts ...client.CallOption) (*grpc.GetRoyaltyReportResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) UnPublishKeyProduct(ctx context.Context, in *grpc.UnPublishKeyProductRequest, opts ...client.CallOption) (*grpc.KeyProductResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) PaymentFormPlatformChanged(
	ctx context.Context,
	in *grpc.PaymentFormUserChangePlatformRequest,
	opts ...client.CallOption,
) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) OrderReceipt(ctx context.Context, in *grpc.OrderReceiptRequest, opts ...client.CallOption) (*grpc.OrderReceiptResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) OrderReceiptRefund(ctx context.Context, in *grpc.OrderReceiptRequest, opts ...client.CallOption) (*grpc.OrderReceiptResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetRecommendedPriceByPriceGroup(ctx context.Context, in *grpc.RecommendedPriceRequest, opts ...client.CallOption) (*grpc.RecommendedPriceResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetRecommendedPriceByConversion(ctx context.Context, in *grpc.RecommendedPriceRequest, opts ...client.CallOption) (*grpc.RecommendedPriceResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) CheckSkuAndKeyProject(ctx context.Context, in *grpc.CheckSkuAndKeyProjectRequest, opts ...client.CallOption) (*grpc.EmptyResponseWithStatus, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetPriceGroupByRegion(ctx context.Context, in *grpc.GetPriceGroupByRegionRequest, opts ...client.CallOption) (*grpc.GetPriceGroupByRegionResponse, error) {
	panic("implement me")
}

func (s *BillingServerErrorMock) GetMerchantUsers(ctx context.Context, in *grpc.GetMerchantUsersRequest, opts ...client.CallOption) (*grpc.GetMerchantUsersResponse, error) {
	panic("implement me")
}
