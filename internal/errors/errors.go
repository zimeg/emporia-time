package errors

const (
	ErrCognitoAuthenticate     errorCode = "err_cognito_authenticate"
	ErrCognitoGenerate         errorCode = "err_cognito_generate"
	ErrCognitoRefresh          errorCode = "err_cognito_refresh"
	ErrCognitoSetup            errorCode = "err_cognito_setup"
	ErrConfigFlag              errorCode = "err_config_flag"
	ErrConfigHelp              errorCode = "err_config_help"
	ErrConfigHome              errorCode = "err_config_home"
	ErrConfigLoad              errorCode = "err_config_load"
	ErrConfigParse             errorCode = "err_config_parse"
	ErrConfigPath              errorCode = "err_config_path"
	ErrConfigRead              errorCode = "err_config_read"
	ErrConfigSave              errorCode = "err_config_save"
	ErrConfigSetup             errorCode = "err_config_setup"
	ErrConfigWrite             errorCode = "err_config_write"
	ErrEmporiaChart            errorCode = "err_emporia_chart"
	ErrEmporiaCheckup          errorCode = "err_emporia_checkup"
	ErrEmporiaComplete         errorCode = "err_emporia_complete"
	ErrEmporiaDevice           errorCode = "err_emporia_device"
	ErrEmporiaDevices          errorCode = "err_emporia_devices"
	ErrEmporiaFormat           errorCode = "err_emporia_format"
	ErrEmporiaMaintenance      errorCode = "err_emporia_maintenance"
	ErrEmporiaMessage          errorCode = "err_emporia_message"
	ErrEmporiaRequest          errorCode = "err_emporia_request"
	ErrEmporiaResponse         errorCode = "err_emporia_response"
	ErrEmporiaResult           errorCode = "err_emporia_result"
	ErrEmporiaStatus           errorCode = "err_emporia_status"
	ErrEmporiaUnplugged        errorCode = "err_emporia_unplugged"
	ErrPromptInput             errorCode = "err_prompt_input"
	ErrPromptSelect            errorCode = "err_prompt_select"
	ErrPromptSelectDescription errorCode = "err_prompt_select_description"
	ErrPromptSelectMissing     errorCode = "err_prompt_select_missing"
	ErrTemplateBuild           errorCode = "err_template_build"
	ErrTemplateFormat          errorCode = "err_template_format"
	ErrTemplateParse           errorCode = "err_template_parse"
	ErrTemplatePrint           errorCode = "err_template_print"
	ErrTimeCommand             errorCode = "err_time_command"
	ErrTimeExecution           errorCode = "err_time_execution"
	ErrTimeParseCode           errorCode = "err_time_parse_code"
	ErrTimeParseReal           errorCode = "err_time_parse_real"
	ErrTimeParseSys            errorCode = "err_time_parse_sys"
	ErrTimeParseUser           errorCode = "err_time_parse_user"
	ErrUnexpectedProblem       errorCode = "err_unexpected_problem" // ErrUnexpectedProblem is a fallback value
	ErrWriteBuffer             errorCode = "err_write_buffer"
	ErrWriteOutput             errorCode = "err_write_output"
)

// New raises a code into an error class
func New(code errorCode) (err Err) {
	defer func() {
		if err.Code != ErrUnexpectedProblem {
			err.Code = code
		}
	}()
	errorMap := map[errorCode]Err{
		ErrCognitoAuthenticate: {
			Message: "failed to initiate authentication",
		},
		ErrCognitoGenerate: {
			Message: "failed to generate tokens",
		},
		ErrCognitoRefresh: {
			Message: "failed to refresh tokens",
		},
		ErrCognitoSetup: {
			Message: "failed to setup provider",
		},
		ErrConfigFlag: {
			Message: "failed to parse config flags",
		},
		ErrConfigHelp: {
			Message: "failed to print config help",
		},
		ErrConfigHome: {
			Message: "failed to find config home",
		},
		ErrConfigLoad: {
			Message: "failed to load config values",
		},
		ErrConfigParse: {
			Message: "failed to parse configurations",
		},
		ErrConfigPath: {
			Message: "failed to make config path",
		},
		ErrConfigRead: {
			Message: "failed to read config file",
		},
		ErrConfigSave: {
			Message: "failed to save config values",
		},
		ErrConfigSetup: {
			Message: "failed to setup configurations",
		},
		ErrConfigWrite: {
			Message: "failed to write config file",
		},
		ErrEmporiaChart: {
			Message: "failed to gather measurement",
		},
		ErrEmporiaCheckup: {
			Message: "failed to check uptime",
		},
		ErrEmporiaComplete: {
			Message: "failed to complete response",
		},
		ErrEmporiaDevice: {
			Message: "failed to find device",
		},
		ErrEmporiaDevices: {
			Message: "failed to gather devices",
		},
		ErrEmporiaFormat: {
			Message: "failed to parse response",
		},
		ErrEmporiaMaintenance: {
			Message: "cannot measure during maintenance",
		},
		ErrEmporiaMessage: {
			Message: "error in get response",
		},
		ErrEmporiaRequest: {
			Message: "failed to make request",
		},
		ErrEmporiaResponse: {
			Message: "failed to get response",
		},
		ErrEmporiaResult: {
			Message: "failed to read result",
		},
		ErrEmporiaStatus: {
			Message: "failed to get status",
		},
		ErrEmporiaUnplugged: {
			Message: "no devices were found",
		},
		ErrPromptInput: {
			Message: "failed to prompt for input value",
		},
		ErrPromptSelect: {
			Message: "failed to prompt for select value",
		},
		ErrPromptSelectDescription: {
			Message: "failed to match select descriptions",
		},
		ErrPromptSelectMissing: {
			Message: "no options to select from",
		},
		ErrTemplateBuild: {
			Message: "failed to build template",
		},
		ErrTemplateFormat: {
			Message: "failed to format template",
		},
		ErrTemplateParse: {
			Message: "failed to parse template",
		},
		ErrTemplatePrint: {
			Message: "failed to print template",
		},
		ErrTimeCommand: {
			Message: "failed to run time command",
		},
		ErrTimeExecution: {
			Message: "failed to execute time command",
		},
		ErrTimeParseCode: {
			Message: "failed to parse exit code",
		},
		ErrTimeParseReal: {
			Message: "failed to parse 'real' time",
		},
		ErrTimeParseSys: {
			Message: "failed to parse 'sys' time",
		},
		ErrTimeParseUser: {
			Message: "failed to parse 'user' time",
		},
		ErrWriteBuffer: {
			Message: "failed to write to buffer",
		},
		ErrWriteOutput: {
			Message: "failed to write to output",
		},
	}
	err, ok := errorMap[code]
	if !ok {
		return Err{
			Code:    ErrUnexpectedProblem,
			Message: string(code),
		}
	}
	return err
}
