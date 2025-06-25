package rabbitmq

const ResearchDataQueue = "research_data"

const BalanceClientsQueue = "balance_clients"

const ChartDataCalculateQueue = "chart_data_calculate"
const ChartDataCalculateRetryInterval = 5000
const ChartDataInsertQueue = "chart_data_insert"
const ChartDataInsertRetryInterval = 5000

const ScheduledOrderQueue = "place_scheduled_order"
const OrderTrackerStatus = "order_tracker_status"
const OrderTrackerEvent = "order_tracker_event"

const DelayFiveSec = 5000
const DelayFiftySec = 50000
const DelayOneHour = 3600 * 1000
const MaxRetriesPublish = 3
const MaxRetriesConsume = 5

const SSLIPNQueue = "ssl_ipn_queue"
const TriggerOrderQueue = "trigger_order_queue"
const MetaInfoQueue = "meta_info_queue"
const UserOrderTrackerQueue = "user_order_tracker_queue"

// DEPOSITS

const WithdrawalsQueue = "backoffice_withdrawals"
const BackofficeDepositsQueue = "backoffice_deposits"
const OMSDepositsQueue = "oms_deposits"

const CreateAccQueue = "create_acc"
const CreateBackOfficeAccQueue = "create_acc_backoffice"
const CreateOmsAccQueue = "create_acc_oms"
const AssignBoQueue = "assign_bo"
const RegisterAccQueue = "register_acc"
const CreateExistingNBLSLBoQueue = "create_existing_nbl_sl_bo"
const BackofficeBalance = "backoffice_balance"
const UserOrderTracker = "user_order_tracker"
const (
	TypeAQueue = "type_a"
	TypeBQueue = "type_b"
	TypeCQueue = "type_c"
	TypeDQueue = "type_d"
	TypeEQueue = "type_e"
	TypeHQueue = "type_h"
	TypeIQueue = "type_i"
	TypeLQueue = "type_l"
	TypeMQueue = "type_m"
	TypeNQueue = "type_n"
	TypePQueue = "type_p"
	TypeQQueue = "type_q"
	TypeRQueue = "type_r"
	TypeSQueue = "type_s"
	TypeTQueue = "type_t"
	TypeUQueue = "type_u"
)
const RejectedReasonTracker = "rejected_reason_tracker"
const FinalizeWithdrawQueue = "finalize_withdraw"

const (
	IntraDayNewOrderQueue    = "intra_day_new_order"
	IntraDayCancelOrderQueue = "intra_day_cancel_order"
	IntraDayModifyOrderQueue = "intra_day_modify_order"
)

const BackofficeRegisterQueue = "backoffice_register"
const BackofficeLimitQueue = "backoffice_limit"
const BackofficeClientQueue = "backoffice_client"
const BackofficeDeactiveQueue = "backoffice_deactive"
const BackofficeInsertOneQueue = "backoffice_insert_one"
