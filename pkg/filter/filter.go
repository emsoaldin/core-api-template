package filter

import (
	"math"
	"strconv"
	"strings"
)

type DeliveryType uint64

const (
	// DeliveryTypeUnknown unknown
	DeliveryTypeUnknown DeliveryType = math.MaxUint64
	// DeliveryTypeElectronic electronic
	DeliveryTypeElectronic DeliveryType = 1
	// DeliveryTypeFlash flash
	DeliveryTypeFlash DeliveryType = 2
	// DeliveryTypePaperless paperless
	DeliveryTypePaperless DeliveryType = 3
	// DeliveryTypeTMET tmet
	DeliveryTypeTMET DeliveryType = 4
	// DeliveryTypeMobileScreencap Mobile Screencap
	DeliveryTypeMobileScreencap DeliveryType = 5
	// DeliveryTypePaperlessCard Paperless Card
	DeliveryTypePaperlessCard DeliveryType = 6
	// DeliveryTypePaperlessWalkin Paperless Walkin
	DeliveryTypePaperlessWalkin DeliveryType = 7
	// DeliveryTypeHard Hard
	DeliveryTypeHard DeliveryType = 8
)

// List of query keys for filtering
var (
	// Setting default value for unit and total price from
	PriceFromDefault = 0
	// Setting default value for unit and total price to
	PriceToDefault = 0

	FilterEventDateFromKey = "event_date_from"

	FilterEventDateToKey = "event_date_to"

	FilterDeliveryDateFromKey = "delivery_date_from"

	FilterDeliveryDateToKey = "delivery_date_to"

	FilterVenueKey = "venue"

	FilterSectionKey = "section"

	FilterUnitPriceFromKey = "min_unit_price"

	FilterUnitPriceToKey = "max_unit_price"

	FilterTotalPriceFromKey = "total_price_from"

	FilterTotalPriceToToKey = "total_price_to"

	FilterDeliveryTypeKey = "delivery_type"

	FilterStatusKey = "status"

	FilterEventNameKey = "event_name"

	FilterOrderNumberKey = "order_number"

	FilterRowKey = "row"

	FilterQuantityKey = "quantity"

	FilterEventIDKey = "event_id"

	FilterFileUploadedKey = "file_uploaded"

	FilterIncludeTagsKey = "include_tags"

	FilterExcludeTagsKey = "exclude_tags"

	FilterInventoryIDKey = "inventory_id"

	FilterSeatTypeKey = "seat_type"

	FilterLowSeatKey = "low_seat"

	FilterHighSeatKey = "high_seat"

	FilterHiddenSeatsKey = "hidden_seats"

	FilterMinQtyKey = "min_qty"

	FilterMaxQtyKey = "max_qty"

	FilterSplitTypeKey = "split_type"

	FilterVendorKey = "vendor_id"

	FilterVendorCompanyKey = "vendor_company"

	FilterVendorFirstNameKey = "vendor_first_name"

	FilterVendorLastNameKey = "vendor_last_name"

	FilterVendorEmailKey = "vendor_email"

	FilterVendorPhoneNumberKey = "vendor_phone_number"

	FilterVendorCityKey = "vendor_city"

	FilterVendorStateKey = "vendor_state"

	FilterExternalNotesKey = "external_notes"

	FilterInternalNotesKey = "internal_notes"

	FilterMinUnitCostKey = "min_unit_cost"

	FilterMaxUnitCostKey = "max_unit_cost"

	FilterCustomerKey = "customer_id"

	FilterCustomerNameKey = "customer_name"

	FilterCustomerEmailKey = "customer_email"

	FilterCustomerPhoneNumberKey = "customer_phone_number"

	FilterCustomerCityKey = "customer_city"

	FilterCustomerStateKey = "customer_state"

	FilterInvoiceIDKey = "invoice_id"

	FilterInvoicePaymentStatusKey = "payment_status"

	FilterInvoiceExternalReferenceKey = "external_reference"

	FilterBarcodesUploadedKey = "barcodes_uploaded"
)

// Filter is a type used to represent the filtering by this params
type Filter struct {
	// Filter from cretain event date
	EventDateFrom string `json:"eventDateFrom"`
	// Filter to cretain event date
	EventDateTo string `json:"eventDateTo"`
	// Filter from cretain delivery date
	DeliveryDateFrom string `json:"deliveryDateFrom"`
	// Filter to cretain delivery date
	DeliveryDateTo string `json:"deliveryDateTo"`
	// Filter by venue name
	Venue string `json:"venue"`
	// Filter by section
	Section string `json:"section"`
	// Filter from certain unit price
	UnitPriceFrom int `json:"unitPriceFrom"`
	// Filter to certain unit price
	UnitPriceTo int `json:"unitPriceTo"`
	// Filter from certain total price
	TotalPriceFrom int `json:"totalPriceFrom"`
	// Filter to certain total price
	TotalPriceTo int `json:"totalPriceTo"`
	// Filter by delivery type
	DeliveryType string `json:"deliveryType"`
	// Filter by status
	Status string `json:"status"`
	// Filter by event name
	EventName string `json:"eventName"`
	// Filter by sale/purchase number
	OrderNumber int `json:"orderNumber"`
	// Filter by row
	Row string `json:"row"`
	// Filter by ticket quantity
	Quantity string `json:"quantity"`
	// Filter by event id
	EventID int `json:"eventId"`
	// Filter by are files uploaded
	FileUploaded bool `json:"fileUploaded"`
	// Filter by tags
	IncludeTags string `json:"includeTags"`
	// Filter by excluded tags
	ExcludeTags string `json:"excludeTags"`
	// Filter by inventory id
	InventoryID int `json:"inventoryId"`
	// Filter by seat type
	SeatType int `json:"seatType"`
	// Filter by low seat
	LowSeat int `json:"lowSeat"`
	// Filter by high seat
	HighSeat int `json:"highSeat"`
	// Filter by hidden seats
	HiddenSeats bool `json:"hiddenSeats"`
	// Filter by min qty
	MinQty int `json:"minQty"`
	// Filter by max qty
	MaxQty int `json:"maxQty"`
	// Filter by split type
	SplitType int `json:"splitType"`
	// Filter by vendor id
	Vendor int `json:"vendorId"`
	// Filter by vendor company
	VendorCompany string `json:"vendorCompany"`
	// Filter by vendor first name
	VendorFirstName string `json:"vendorFirstName"`
	// Filter by vendor last name
	VendorLastName string `json:"vendorLastName"`
	// Filter by vendor email
	VendorEmail string `json:"vendorEmail"`
	// Filter by vendor phone number
	VendorPhoneNumber string `json:"vendorPhoneNumber"`
	// Filter by vendor city
	VendorCity string `json:"vendorCity"`
	// Filter by vendor state
	VendorState string `json:"vendorState"`
	// Filter by public notes
	PublicNotes string `json:"publicNotes"`
	// Filter by external notes
	InternalNotes string `json:"internalNotes"`
	// Filter by min unit cost
	MinUnitCost int `json:"minUnitCost"`
	// Filter by max unit cost
	MaxUnitCost int `json:"maxUnitCost"`
	// Filter by customer
	Customer int `json:"customerId"`
	// Filter by customer name
	CustomerName string `json:"customerFirstName"`
	// Filter by customer email
	CustomerEmail string `json:"customerEmail"`
	// Filter by customer phone number
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
	// Filter by customer city
	CustomerCity string `json:"customerCity"`
	// Filter by customer state
	CustomerState string `json:"customerState"`
	// Filter by invoice id
	InvoiceID int `json:"invoiceId"`
	// Filter by payment status
	InvoicePaymentStatus int `json:"invoicePaymentId"`
	// Filter by external reference
	InvoiceExternalReference string `json:"invoiceExternalReference"`
	// Filter by barcodes uploaded
	BarcodesUploaded bool `json:"barcodesUploaded"`
}

// PaginationParams is a parameters provider interface to get the pagination params from
type FilterParams interface {
	Get(key string) string
}

// NewPaginator returns a new `Paginator` value with the appropriate
// defaults set.
func NewFilter(eventDateFrom string, eventDateTo string,
	deliveryDateFrom string,
	deliveryDateTo string,
	venue string,
	section string,
	unitPriceFrom int,
	unitPriceTo int,
	totalPriceFrom int,
	totalPriceTo int,
	deliveryType string,
	statusIDs string,
	eventName string,
	orderNumber int,
	row string,
	quantity string,
	eventId int,
	includeTags string,
	excludeTags string,
	inventoryID int,
	seatType int,
	lowSeat int,
	highSeat int,
	minQty int,
	maxQty int,
	splitType int,
	vendor int,
	internalNotes string,
	publicNotes string,
	minUnitCost int,
	maxUnitCost int,
	fileUploaded bool,
	hiddenSeats bool,
	customer int,
	vendorCompany string,
	vendorFirstName string,
	vendorLastName string,
	vendorEmail string,
	vendorPhoneNumber string,
	vendorCity string,
	vendorState string,
	customerName string,
	customerEmail string,
	customerPhoneNumber string,
	customerCity string,
	customerState string,
	invoiceId int,
	invoicePaymentStatus int,
	invoiceExternalReference string,
	barcodesUploaded bool) *Filter {

	f := &Filter{EventDateFrom: eventDateFrom,
		EventDateTo:              eventDateTo,
		DeliveryDateFrom:         deliveryDateFrom,
		DeliveryDateTo:           deliveryDateTo,
		Venue:                    venue,
		Section:                  section,
		UnitPriceFrom:            unitPriceFrom,
		UnitPriceTo:              unitPriceTo,
		TotalPriceFrom:           totalPriceFrom,
		TotalPriceTo:             totalPriceTo,
		DeliveryType:             deliveryType,
		Status:                   statusIDs,
		EventName:                eventName,
		OrderNumber:              orderNumber,
		Row:                      row,
		Quantity:                 quantity,
		EventID:                  eventId,
		FileUploaded:             fileUploaded,
		IncludeTags:              includeTags,
		ExcludeTags:              excludeTags,
		InventoryID:              inventoryID,
		SeatType:                 seatType,
		LowSeat:                  lowSeat,
		HighSeat:                 highSeat,
		HiddenSeats:              hiddenSeats,
		MinQty:                   minQty,
		MaxQty:                   maxQty,
		SplitType:                splitType,
		Vendor:                   vendor,
		VendorCompany:            vendorCompany,
		VendorFirstName:          vendorFirstName,
		VendorLastName:           vendorLastName,
		VendorEmail:              vendorEmail,
		VendorPhoneNumber:        vendorPhoneNumber,
		VendorCity:               vendorCity,
		VendorState:              vendorState,
		InternalNotes:            internalNotes,
		PublicNotes:              publicNotes,
		MinUnitCost:              minUnitCost,
		MaxUnitCost:              maxUnitCost,
		Customer:                 customer,
		CustomerName:             customerName,
		CustomerEmail:            customerEmail,
		CustomerPhoneNumber:      customerPhoneNumber,
		CustomerCity:             customerCity,
		CustomerState:            customerState,
		InvoiceID:                invoiceId,
		InvoicePaymentStatus:     invoicePaymentStatus,
		InvoiceExternalReference: invoiceExternalReference,
		BarcodesUploaded:         barcodesUploaded}
	return f
}

// NewPaginatorFromParams takes an interface of type `PaginationParams`,
// the `url.Values` type works great with this interface, and returns
// a new `Paginator` based on the params or `PaginatorPageKey` and
// `PaginatorPerPageKey`. Defaults are `1` for the page and
// PaginatorPerPageDefault for the per page value.
func NewFilterFromParams(params FilterParams) *Filter {
	eventDateFrom := ""
	if edf := params.Get(FilterEventDateFromKey); edf != "" {
		eventDateFrom = strings.TrimSpace(edf)
	}

	eventDateTo := ""
	if edt := params.Get(FilterEventDateToKey); edt != "" {
		eventDateTo = strings.TrimSpace(edt)
	}

	deliveryDateFrom := ""
	if ddf := params.Get(FilterDeliveryDateFromKey); ddf != "" {
		deliveryDateFrom = strings.TrimSpace(ddf)
	}

	deliveryDateTo := ""
	if ddt := params.Get(FilterDeliveryDateToKey); ddt != "" {
		deliveryDateTo = strings.TrimSpace(ddt)
	}

	venue := ""
	if v := params.Get(FilterVenueKey); v != "" {
		venue = strings.TrimSpace(v)
	}

	section := ""
	if s := params.Get(FilterSectionKey); s != "" {
		section = strings.TrimSpace(s)
	}

	unitPriceFrom := strconv.Itoa(PriceFromDefault)
	if upf := params.Get(FilterUnitPriceFromKey); upf != "" {
		unitPriceFrom = upf
	}

	upf, err := strconv.Atoi(unitPriceFrom)
	if err != nil {
		upf = PriceFromDefault
	}

	unitPriceTo := strconv.Itoa(PriceToDefault)
	if upt := params.Get(FilterUnitPriceToKey); upt != "" {
		unitPriceTo = upt
	}

	upt, err := strconv.Atoi(unitPriceTo)
	if err != nil {
		upt = PriceToDefault
	}

	totalPriceFrom := strconv.Itoa(PriceFromDefault)
	if tpf := params.Get(FilterTotalPriceFromKey); tpf != "" {
		totalPriceFrom = tpf
	}

	tpf, err := strconv.Atoi(totalPriceFrom)
	if err != nil {
		tpf = PriceFromDefault
	}

	totalPriceTo := strconv.Itoa(PriceToDefault)
	if tpt := params.Get(FilterTotalPriceToToKey); tpt != "" {
		totalPriceTo = tpt
	}

	tpt, err := strconv.Atoi(totalPriceTo)
	if err != nil {
		tpt = PriceToDefault
	}

	deliveryType := ""
	if dt := params.Get(FilterDeliveryTypeKey); dt != "" {
		deliveryType = strings.TrimSpace(dt)
	}

	statusIDs := ""
	if st := params.Get(FilterStatusKey); st != "" {
		statusIDs = strings.TrimSpace(st)
	}

	eventName := ""
	if en := params.Get(FilterEventNameKey); en != "" {
		eventName = strings.TrimSpace(en)
	}

	orderNumber := strconv.Itoa(0)
	if on := params.Get(FilterOrderNumberKey); on != "" {
		orderNumber = on
	}

	on, err := strconv.Atoi(orderNumber)
	if err != nil {
		on = 0
	}

	row := ""
	if ro := params.Get(FilterRowKey); ro != "" {
		row = ro
	}

	quantity := ""
	if qu := params.Get(FilterQuantityKey); qu != "" {
		quantity = qu
	}

	eventId := strconv.Itoa(PriceToDefault)
	if eid := params.Get(FilterEventIDKey); eid != "" {
		eventId = eid
	}

	eid, err := strconv.Atoi(eventId)
	if err != nil {
		eid = PriceToDefault
	}

	includeTags := ""
	if it := params.Get(FilterIncludeTagsKey); it != "" {
		includeTags = it
	}

	excludeTags := ""
	if et := params.Get(FilterExcludeTagsKey); et != "" {
		excludeTags = et
	}

	inventoryId := strconv.Itoa(PriceToDefault)
	if iid := params.Get(FilterInventoryIDKey); iid != "" {
		inventoryId = iid
	}

	iid, err := strconv.Atoi(inventoryId)
	if err != nil {
		iid = PriceToDefault
	}

	seatType := strconv.Itoa(PriceToDefault)
	if st := params.Get(FilterSeatTypeKey); st != "" {
		seatType = st
	}

	st, err := strconv.Atoi(seatType)
	if err != nil {
		st = PriceToDefault
	}

	lowSeat := strconv.Itoa(PriceToDefault)
	if ls := params.Get(FilterLowSeatKey); ls != "" {
		lowSeat = ls
	}

	ls, err := strconv.Atoi(lowSeat)
	if err != nil {
		ls = PriceToDefault
	}

	highSeat := strconv.Itoa(PriceToDefault)
	if hs := params.Get(FilterHighSeatKey); hs != "" {
		highSeat = hs
	}

	hs, err := strconv.Atoi(highSeat)
	if err != nil {
		hs = PriceToDefault
	}

	minQty := strconv.Itoa(PriceToDefault)
	if minq := params.Get(FilterMinQtyKey); minq != "" {
		minQty = minq
	}

	minq, err := strconv.Atoi(minQty)
	if err != nil {
		minq = PriceToDefault
	}

	maxQty := strconv.Itoa(PriceToDefault)
	if maxq := params.Get(FilterMaxQtyKey); maxq != "" {
		maxQty = maxq
	}

	maxq, err := strconv.Atoi(maxQty)
	if err != nil {
		maxq = PriceToDefault
	}

	splitType := strconv.Itoa(PriceToDefault)
	if spt := params.Get(FilterSplitTypeKey); spt != "" {
		splitType = spt
	}

	spt, err := strconv.Atoi(splitType)
	if err != nil {
		spt = PriceToDefault
	}

	vendor := strconv.Itoa(PriceToDefault)
	if vid := params.Get(FilterVendorKey); vid != "" {
		vendor = vid
	}

	vid, err := strconv.Atoi(vendor)
	if err != nil {
		vid = PriceToDefault
	}

	internalNotes := ""
	if in := params.Get(FilterInternalNotesKey); in != "" {
		internalNotes = in
	}

	externalNotes := ""
	if en := params.Get(FilterExternalNotesKey); en != "" {
		externalNotes = en
	}

	minCost := strconv.Itoa(PriceToDefault)
	if minc := params.Get(FilterMinUnitCostKey); minc != "" {
		minCost = minc
	}

	minc, err := strconv.Atoi(minCost)
	if err != nil {
		minc = PriceToDefault
	}

	maxCost := strconv.Itoa(PriceToDefault)
	if maxc := params.Get(FilterMaxUnitCostKey); maxc != "" {
		maxCost = maxc
	}

	maxc, err := strconv.Atoi(maxCost)
	if err != nil {
		maxc = PriceToDefault
	}

	fileUploaded := false
	if fp := params.Get(FilterFileUploadedKey); fp != "" {
		fileUploaded = true
	}

	hiddenSeats := false
	if fp := params.Get(FilterHiddenSeatsKey); fp != "" {
		hiddenSeats = true
	}

	customer := strconv.Itoa(PriceToDefault)
	if cust := params.Get(FilterCustomerKey); cust != "" {
		customer = cust
	}

	cust, err := strconv.Atoi(customer)
	if err != nil {
		cust = PriceToDefault
	}

	vendorCompany := ""
	if vc := params.Get(FilterVendorCompanyKey); vc != "" {
		vendorCompany = vc
	}

	vendorFirstName := ""
	if vfn := params.Get(FilterVendorFirstNameKey); vfn != "" {
		vendorFirstName = vfn
	}

	vendorLastName := ""
	if vln := params.Get(FilterVendorLastNameKey); vln != "" {
		vendorLastName = vln
	}

	vendorEmail := ""
	if ve := params.Get(FilterVendorEmailKey); ve != "" {
		vendorEmail = ve
	}

	vendorPhoneNumber := ""
	if vpn := params.Get(FilterVendorPhoneNumberKey); vpn != "" {
		vendorPhoneNumber = vpn
	}

	vendorCity := ""
	if vcit := params.Get(FilterVendorCityKey); vcit != "" {
		vendorCity = vcit
	}

	vendorState := ""
	if vs := params.Get(FilterVendorStateKey); vs != "" {
		vendorState = vs
	}

	customerName := ""
	if vln := params.Get(FilterCustomerNameKey); vln != "" {
		customerName = vln
	}

	customerEmail := ""
	if ve := params.Get(FilterCustomerEmailKey); ve != "" {
		customerEmail = ve
	}

	customerPhoneNumber := ""
	if vpn := params.Get(FilterCustomerPhoneNumberKey); vpn != "" {
		customerPhoneNumber = vpn
	}

	customerCity := ""
	if vcit := params.Get(FilterCustomerCityKey); vcit != "" {
		customerCity = vcit
	}

	customerState := ""
	if vs := params.Get(FilterCustomerStateKey); vs != "" {
		customerState = vs
	}

	inid := strconv.Itoa(PriceToDefault)
	if invid := params.Get(FilterInvoiceIDKey); invid != "" {
		inid = invid
	}

	invid, err := strconv.Atoi(inid)
	if err != nil {
		invid = PriceToDefault
	}

	ipid := strconv.Itoa(PriceToDefault)
	if invpid := params.Get(FilterInvoicePaymentStatusKey); invpid != "" {
		ipid = invpid
	}

	invpid, err := strconv.Atoi(ipid)
	if err != nil {
		invpid = PriceToDefault
	}

	invoiceExternalReference := ""
	if ier := params.Get(FilterInvoiceExternalReferenceKey); ier != "" {
		invoiceExternalReference = ier
	}

	barcodesUploaded := false
	if fp := params.Get(FilterBarcodesUploadedKey); fp != "" {
		barcodesUploaded = true
	}

	return NewFilter(eventDateFrom,
		eventDateTo,
		deliveryDateFrom,
		deliveryDateTo,
		venue,
		section,
		upf,
		upt,
		tpf,
		tpt,
		deliveryType,
		statusIDs,
		eventName,
		on,
		row,
		quantity,
		eid,
		includeTags,
		excludeTags,
		iid,
		st,
		ls,
		hs,
		minq,
		maxq,
		spt,
		vid,
		internalNotes,
		externalNotes,
		minc,
		maxc,
		fileUploaded,
		hiddenSeats,
		cust,
		vendorCompany,
		vendorFirstName,
		vendorLastName,
		vendorEmail,
		vendorPhoneNumber,
		vendorCity,
		vendorState,
		customerName,
		customerEmail,
		customerPhoneNumber,
		customerCity,
		customerState,
		invid,
		invpid,
		invoiceExternalReference,
		barcodesUploaded)
}
