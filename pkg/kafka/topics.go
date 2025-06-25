package kafka

const InAppNotificationTopic = "in_app_notification"
const InAppNotificationGroup = "in_app_notification_grp_1"
const CreateBoTopic = "createBO"
const CreateBoGroup = "create_BO_grp_1"
const GlobalNotificationTopic = "global_notification"
const GlobalNotificationGroup = "global_notification_grp_1"
const EmailPusherTopic = "email-pusher"
const EmailPusherGroup = "email-pusher-grp-1"
const BankOtpVerifierTopic = "bank-otp-verifier"
const BankOtpVerifierGroup = "bank-otp-verifier-grp-1"
const SMSPusherTopic = "sms-pusher"
const SMSPusherGroup = "sms-pusher-grp-1"
const DSEIndexTopic = "dse-index"
const DSEIndexGroup = "dse-index-grp-1"
const DSEIndexPartitions = 3
const PushNotificationGroup = "push_notification"
const PushNotificationTopic = "push_notification_grp_1"
const ResearchFileUploadTopic = "research_file_upload"
const ResearchFileUploadGroup = "research_file_upload_grp_1"
const TickerStatusTopic = "ticker_status"
const TickerStatusGroup = "ticker_status_grp_1"
const TickerStatusPartitions = 3
const MarketDepthTopic = "market_depth"
const MarketDepthGroup = "market_depth_grp_1"
const MarketDepthPartitions = 3
const DailyMoversTopic = "daily_movers"
const DailyMoversGroup = "daily_movers_grp_1"
const DailyMoversPartitions = 3
const MarketDataParseTopic = "market_data_parse"
const MarketDataParseGroup = "market_data_parse_grp_1"
const MarketDataParsePartitions = 5
const TradingHistoryTopic = "trading_history"
const TradingHistoryGroup = "trading_history_grp_1"
const TradingHistoryPartitions = 3
const TickerAlertTopic = "ticker_alert"
const TickerAlertGroup = "ticker_alert_grp_1"
const ChartTopic = "chart_topic"
const ChartGroup = "chart_topic_grp_1"
const ChartPartitions = 3
const AnnouncementTopic = "announcement_topic"
const AnnouncementGroup = "announcement_topic_grp_1"
const ForcedSignOutTopic = "forced_sign_out"
const ForcedSignOutGroup = "forced_sign_out_grp_1"
const OrderStatusTopic = "order_status"
const OrderStatusTopicGroup = "order_status_grp_1"
const UserInfoTopic = "user_info"
const UserInfoGroup = "user_info_grp_1"
const AddOrUpdateFCMSignatureTopic = "add_or_update_fcm_signature_topic"
const AddOrUpdateFCMSignatureGroup = "add_or_update_fcm_signature_grp_1"
const WebhookOrderStatusTopic = "webhook_order_status"
const WebhookOrderStatusGroup = "webhook_order_status_grp_1"
const PushNotificationDeviceTopic = "push_notification_device"
const PushNotificationDeviceGroup = "push_notification_device_grp_1"

const ItchHistory = "itch_history"
const ItchHistoryGroup = "itch_history_grp"

const ItchSequenceData = "itch_seq_msg"
const ItchSequenceDataGroup = "itch_seq_group"

const TypeATopic = "type_a_topic"
const TypeAGroup = "type_a_group"
const TypeAGroupForOrderState = "type_a_group_order_state"
const TypeAGroupForTimeNSales = "type_a_group_time_n_sales"
const TypeAGroupForActiveTickerInfo = "type_a_group_active_ticker_info"

const TypeFTopic = "type_f_topic"
const TypeFGroup = "type_f_group"

const TypeBTopic = "type_b_topic"
const TypeBGroup = "type_b_group"

const TypeCTopic = "type_c_topic"
const TypeCGroup = "type_c_group"
const TypeCGroupForOrderState = "type_c_group_order_state"

const TypeDTopic = "type_d_topic"
const TypeDGroup = "type_d_group"
const TypeDGroupForOrderState = "type_d_group_order_state"

const TypeETopic = "type_e_topic"
const TypeEGroup = "type_e_group"
const TypeEGroupForOrderState = "type_e_group_order_state"
const TypeEGroupForTimeNSales = "type_e_group_time_n_sales"
const TypeEGroupForActiveTickerInfo = "type_e_group_active_ticker_info"

const TypeHTopic = "type_h_topic"
const TypeHGroup = "type_h_group"
const TypeHGroupForPortfolioTickerStatus = "type_h_group_for_portfolio_ticker_status"

const TypeITopic = "type_i_topic"
const TypeIGroup = "type_i_group"

const TypeLTopic = "type_l_topic"
const TypeLGroup = "type_l_group"

const TypeMTopic = "type_m_topic"
const TypeMGroup = "type_m_group"

const TypeNTopic = "type_n_topic"
const TypeNGroup = "type_n_group"

const TypePTopic = "type_p_topic"
const TypePGroup = "type_p_group"

const TypeQTopic = "type_q_topic"
const TypeQGroup = "type_q_group"
const TypeQGroupForActiveTickerInfo = "type_q_group_active_ticker_info"

const TypeRTopic = "type_r_topic"
const TypeRGroup = "type_r_group"
const TypeRGroupForActiveTickerInfo = "type_r_group_active_ticker_info"
const TypeRGroupForPortfolioTickerInfo = "type_r_group_for_portfolio_ticker_info"

const TypeSTopic = "type_s_topic"
const TypeSGroup = "type_s_group"

const TypeTTopic = "type_t_topic"
const TypeTGroup = "type_t_group"

const TypeUTopic = "type_u_topic"
const TypeUGroup = "type_u_group"
const TypeUGroupForOrderState = "type_u_group_order_state"
const TypeUGroupForTimeNSales = "type_u_group_time_n_sales"

// Add more topics and groups for types V to Z if needed
const TypeVTopic = "type_v_topic"
const TypeVGroup = "type_v_group"

const TypeWTopic = "type_w_topic"
const TypeWGroup = "type_w_group"

const TypeXTopic = "type_x_topic"
const TypeXGroup = "type_x_group"

const TypeYTopic = "type_y_topic"
const TypeYGroup = "type_y_group"

const TypeZTopic = "type_z_topic"
const TypeZGroup = "type_z_group"

const (
	SimulatorData                             = "simulator_data"
	SimulatorDataGroup                        = "simulator_data_group"
	ItchSequenceDataStorageGroup              = "itch_sequence_data_storage_group"
	ItchSequenceDataPlainTextGroup            = "itch_sequence_data_plain_text_group"
	ItchOrderProgress                         = "itch_order_progress"
	ItchOrderProgressGroupForOrderState       = "itch_order_progress_group_for_order_state"
	ItchOrderProgressGroupForActiveTickerInfo = "itch_order_progress_group_for_active_ticker_info"
	SocketActiveTickerInfo                    = "socket_active_ticker_info"
	SocketActiveTickerInfoGroup               = "socket_active_ticker_info_group"
	SocketActiveTickerInfoPortfolioGroup      = "socket_active_ticker_info_portfolio_group"
	SocketActiveTickerInfoTrekGroup           = "socket_active_ticker_info_trek_group"
	ItchIndex                                 = "itch_index"
	ItchIndexGroup                            = "itch_index_group"
)

const (
	PositionValuation            = "position_valuation"
	PositionValuationSocketGroup = "position_valuation_socket_grp"
)

const (
	ExecutionReport                     = "execution_report"
	ExecutionReportGroup                = "execution_report_group"
	ExecutionReportTradeCaptureGroup    = "execution_report_trade_capture_group"
	ExecutionReportAccountsGroup        = "execution_report_accounts_group"
	ExecutionReportInAppGroup           = "execution_report_inapp_group"
	CancelRejectReport                  = "cancel_reject_report"
	CancelRejectReportGroup             = "cancel_reject_report_group"
	CancelRejectReportTradeCaptureGroup = "cancel_reject_report_trade_capture_group"
)

const ErrorTopic = "error_topic"
const ErrorTopicSlackGroup = "error_topic_slack_group"

const UpdateCurrentPriceTopic = "update_current_price_topic"
const UpdateCurrentPriceGroup = "update_current_price_group"
