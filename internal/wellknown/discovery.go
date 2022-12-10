package wellknown

import (
	"encoding/json"
	"errors"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
	"github.com/absurdlab/tigerd/internal/should"
	"github.com/absurdlab/tigerd/internal/spec"
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/samber/lo"
	"io"
	"os"
	"strings"
)

var (
	// ErrDiscovery is the root error returned when Discovery cannot be parsed from JSON or contains invalid values.
	ErrDiscovery = errors.New("discovery is invalid")
)

// Discovery models the OpenID Connect configuration metadata.
type Discovery struct {
	Issuer                                     string                      `json:"issuer,omitempty"`
	AuthorizationEndpoint                      string                      `json:"authorization_endpoint,omitempty"`
	ResumeAuthorizationEndpoint                string                      `json:"resume_authorization_endpoint,omitempty"`
	TokenEndpoint                              string                      `json:"token_endpoint,omitempty"`
	UserInfoEndpoint                           string                      `json:"userinfo_endpoint,omitempty"`
	JSONWebKeySetURI                           string                      `json:"jwks_uri,omitempty"`
	RegistrationEndpoint                       string                      `json:"registration_endpoint,omitempty"`
	ScopesSupported                            []string                    `json:"scopes_supported,omitempty"`
	ResponseTypesSupported                     []spec.ResponseTypeSet      `json:"response_types_supported,omitempty"`
	ResponseModesSupported                     []spec.ResponseMode         `json:"response_modes_supported,omitempty"`
	GrantTypesSupported                        []spec.GrantType            `json:"grant_types_supported,omitempty"`
	AcrValuesSupported                         []string                    `json:"acr_values_supported,omitempty"`
	SubjectTypesSupported                      []spec.SubjectType          `json:"subject_types_supported,omitempty"`
	IdTokenSigningAlgValuesSupported           []spec.SignatureAlgorithm   `json:"id_token_signing_alg_values_supported,omitempty"`
	IdTokenEncryptionAlgValuesSupported        []spec.EncryptionAlgorithm  `json:"id_token_encryption_alg_values_supported,omitempty"`
	IdTokenEncryptionEncValuesSupported        []spec.EncryptionEncoding   `json:"id_token_encryption_enc_values_supported,omitempty"`
	UserInfoSigningAlgValuesSupported          []spec.SignatureAlgorithm   `json:"userinfo_signing_alg_values_supported,omitempty"`
	UserInfoEncryptionAlgValuesSupported       []spec.EncryptionAlgorithm  `json:"userinfo_encryption_alg_values_supported,omitempty"`
	UserInfoEncryptionEncValuesSupported       []spec.EncryptionEncoding   `json:"userinfo_encryption_enc_values_supported,omitempty"`
	RequestObjectSigningAlgValuesSupported     []spec.SignatureAlgorithm   `json:"request_object_signing_alg_values_supported,omitempty"`
	RequestObjectEncryptionAlgValuesSupported  []spec.EncryptionAlgorithm  `json:"request_object_encryption_alg_values_supported,omitempty"`
	RequestObjectEncryptionEncValuesSupported  []spec.EncryptionEncoding   `json:"request_object_encryption_enc_values_supported,omitempty"`
	TokenEndpointAuthMethodsSupported          []spec.AuthenticationMethod `json:"token_endpoint_auth_methods_supported,omitempty"`
	TokenEndpointAuthSigningAlgValuesSupported []spec.SignatureAlgorithm   `json:"token_endpoint_auth_signing_alg_values_supported,omitempty"`
	DisplayValuesSupported                     []spec.Display              `json:"display_values_supported,omitempty"`
	ClaimTypesSupported                        []spec.ClaimType            `json:"claim_types_supported,omitempty"`
	ClaimsSupported                            []string                    `json:"claims_supported,omitempty"`
	ServiceDocumentation                       string                      `json:"service_documentation,omitempty"`
	UILocalesSupported                         []string                    `json:"ui_locales_supported,omitempty"`
	ClaimsParameterSupported                   bool                        `json:"claims_parameter_supported,omitempty"`
	RequestParameterSupported                  bool                        `json:"request_parameter_supported,omitempty"`
	RequestURIParameterSupported               bool                        `json:"request_uri_parameter_supported,omitempty"`
	RequireRequestURIRegistration              bool                        `json:"require_request_uri_registration,omitempty"`
	OPPolicyURI                                string                      `json:"op_policy_uri,omitempty"`
	OPTermsOfServiceURI                        string                      `json:"op_tos_uri,omitempty"`
}

// Apply runs the supplied functions on this Discovery, and potentially modifies this Discovery.
func (d *Discovery) Apply(fn ...func(discovery *Discovery)) {
	for _, it := range fn {
		it(d)
	}
}

func (d *Discovery) onlyImplicitFlow() bool {
	return lo.None(d.GrantTypesSupported, []spec.GrantType{
		spec.GrantTypeAuthorizationCode,
		spec.GrantTypeClientCredentials,
		spec.GrantTypeRefreshToken,
	})
}

// Validate performs validation on the Discovery and returns an error if exists violation. The returned error will
// be a ErrDiscovery.
func (d *Discovery) Validate() error {
	err := v.Errors{
		"issuer": v.Validate(d.Issuer,
			v.Required,
			is.URL,
			should.URL().Https().NoQuery().NoFragment(),
		),
		"authorization_endpoint": v.Validate(d.AuthorizationEndpoint,
			v.Required,
			is.URL,
			should.URL().Http().Https().NoFragment(),
		),
		"token_endpoint": v.Validate(d.TokenEndpoint,
			v.When(!d.onlyImplicitFlow(),
				v.Required,
				is.URL,
				should.URL().Http().Https().NoFragment(),
			),
		),
		"userinfo_endpoint": v.Validate(d.UserInfoEndpoint,
			v.Required,
			is.URL,
			should.URL().Https().NoFragment(),
		),
		"jwks_uri": v.Validate(d.JSONWebKeySetURI,
			v.Required,
			is.URL,
			should.URL().Https().NoFragment(),
		),
		"registration_endpoint": v.Validate(d.RegistrationEndpoint,
			is.URL,
			should.URL().Https().NoFragment(),
		),
		"scopes_supported": v.Validate(d.ScopesSupported,
			v.Required,
			should.Contain(spec.ScopeOpenID).Error("should contain openid"),
		),
		"response_types_supported": v.Validate(d.ResponseTypesSupported,
			v.Required,
			should.Contain(spec.ResponseTypeCode.ToSet(), spec.ResponseTypeIDToken.ToSet()).
				Error("should contain code and id_token"),
		),
		"grant_types_supported": v.Validate(d.GrantTypesSupported,
			v.Required,
			should.Contain(spec.GrantTypeAuthorizationCode, spec.GrantTypeImplicit).
				Error("should contain authorization_code and implicit"),
		),
		"response_modes_supported": v.Validate(d.ResponseModesSupported,
			v.Required,
		),
		"subject_types_supported": v.Validate(d.SubjectTypesSupported,
			v.Required,
		),
		"id_token_signing_alg_values_supported": v.Validate(d.IdTokenSigningAlgValuesSupported,
			v.Required,
			should.Contain(spec.RS256).Error("should contain RS256"),
		),
		"request_object_signing_alg_values_supported": v.Validate(d.RequestObjectSigningAlgValuesSupported,
			v.Required,
			should.Contain(spec.RS256, spec.NoSignature).Error("should contain RS256 and none"),
		),
		"token_endpoint_auth_signing_alg_values_supported": v.Validate(d.TokenEndpointAuthSigningAlgValuesSupported,
			v.Required,
			should.Contain(spec.RS256).Error("should contain RS256"),
		),
		"service_documentation": v.Validate(d.ServiceDocumentation,
			is.URL,
			should.URL().Http().Https(),
		),
		"op_policy_uri": v.Validate(d.OPPolicyURI,
			is.URL,
			should.URL().Http().Https(),
		),
		"op_tos_uri": v.Validate(d.OPTermsOfServiceURI,
			is.URL,
			should.URL().Http().Https(),
		),
	}.Filter()

	if err != nil {
		return fault.Wrap(ErrDiscovery,
			ftag.With(spec.ErrKindInvalidRequest),
			fmsg.WithDesc(err.Error(), err.Error()),
		)
	}

	return nil
}

// DiscoveryProperties is the configuration properties for reading Discovery. Discovery, in its json format, can be
// read from a File, or read from an Inline string. The File option, if specified, precedes the Inline option.
type DiscoveryProperties struct {
	File           string `json:"file" yaml:"file"`
	Inline         string `json:"inline" yaml:"inline"`
	SkipValidation bool   `json:"skipValidation" yaml:"skipValidation"`

	// PreValidateHook is a test-facing hook to modify the sourced Discovery before it is validated. Test cases can
	// modify a correct version of Discovery to create errors, in order to test validation logic.
	PreValidateHook func(d *Discovery) `json:"-" yaml:"-"`
}

// NewDiscovery reads a Discovery by means specified in DiscoveryProperties. By default, validation is performed on
// the sourced Discovery. Setting DiscoveryProperties.SkipValidation to true will skip validation.
func NewDiscovery(props *DiscoveryProperties) (*Discovery, error) {
	var (
		reader io.Reader
		err    error
	)
	switch {
	case len(props.File) > 0:
		reader, err = os.Open(props.File)
		if err != nil {
			err = fault.Wrap(ErrDiscovery,
				ftag.With(spec.ErrKindInvalidRequest),
				fmsg.WithDesc(err.Error(), "Failed to open discovery file."),
			)
		}
	case len(props.Inline) > 0:
		reader, err = strings.NewReader(props.Inline), nil
	default:
		err = fault.Wrap(ErrDiscovery,
			ftag.With(spec.ErrKindInvalidRequest),
			fmsg.WithDesc(err.Error(), "Expect either file or inline option for discovery."),
		)
	}
	if err != nil {
		return nil, err
	}

	discovery := &Discovery{
		ResponseModesSupported:            []spec.ResponseMode{spec.ResponseModeQuery, spec.ResponseModeFragment},
		GrantTypesSupported:               []spec.GrantType{spec.GrantTypeAuthorizationCode, spec.GrantTypeImplicit},
		TokenEndpointAuthMethodsSupported: []spec.AuthenticationMethod{spec.ClientSecretBasic},
		ClaimTypesSupported:               []spec.ClaimType{spec.ClaimTypeNormal},
		RequestURIParameterSupported:      true,
	} // default values

	if err = json.NewDecoder(reader).Decode(&discovery); err != nil {
		return nil, fault.Wrap(ErrDiscovery,
			ftag.With(spec.ErrKindInvalidRequest),
			fmsg.WithDesc("discovery json", err.Error()),
		)
	}

	if props.PreValidateHook != nil {
		props.PreValidateHook(discovery)
	}

	if !props.SkipValidation {
		if err = discovery.Validate(); err != nil {
			return nil, fault.Wrap(err)
		}
	}

	return discovery, nil
}
