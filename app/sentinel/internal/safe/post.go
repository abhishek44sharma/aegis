/*
 * .-'_.---._'-.
 * ||####|(__)||   Protect your secrets, protect your business.
 *   \\()|##//       Secure your sensitive data with Aegis.
 *    \\ |#//                    <aegis.ist>
 *     .\_/.
 */

package safe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	data "github.com/shieldworks/aegis/core/entity/data/v1"
	reqres "github.com/shieldworks/aegis/core/entity/reqres/safe/v1"
	"github.com/shieldworks/aegis/core/env"
	"github.com/shieldworks/aegis/core/validation"
	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"io"
	"log"
	"net/http"
	"net/url"
)

func createAuthorizer() tlsconfig.Authorizer {
	return tlsconfig.AdaptMatcher(func(id spiffeid.ID) error {
		if validation.IsSafe(id.String()) {
			return nil
		}

		return errors.New("Post: I don’t know you, and it’s crazy: '" +
			id.String() + "'",
		)
	})
}

func decideBackingStore(backingStore string) data.BackingStore {
	switch data.BackingStore(backingStore) {
	case data.File:
		return data.File
	case data.Memory:
		return data.Memory
	default:
		return env.SafeBackingStore()
	}
}

func decideSecretFormat(format string) data.SecretFormat {
	switch data.SecretFormat(format) {
	case data.Json:
		return data.Json
	case data.Yaml:
		return data.Yaml
	default:
		return data.Json
	}
}

func newSecretUpsertRequest(workloadId, secret, namespace, backingStore string,
	useKubernetes bool, template string, format string, encrypt, appendSecret bool,
) reqres.SecretUpsertRequest {
	bs := decideBackingStore(backingStore)
	f := decideSecretFormat(format)

	return reqres.SecretUpsertRequest{
		WorkloadId:    workloadId,
		BackingStore:  bs,
		Namespace:     namespace,
		UseKubernetes: useKubernetes,
		Template:      template,
		Format:        f,
		Encrypt:       encrypt,
		AppendValue:   appendSecret,
		Value:         secret,
	}
}

func respond(r *http.Response) {
	if r == nil {
		return
	}

	defer func(b io.ReadCloser) {
		if b == nil {
			return
		}
		err := b.Close()
		if err != nil {
			log.Println("Post: Problem closing request body.", err.Error())
		}
	}(r.Body)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Post: Unable to read the response body from Aegis Safe.", err.Error())
		fmt.Println("")
		return
	}

	fmt.Println("")
	fmt.Println(string(body))
	fmt.Println("")
}

func printEndpointError(err error) {
	fmt.Println("Post: I am having problem generating Aegis Safe "+
		"secrets api endpoint URL.", err.Error())
	fmt.Println("")
}

func printPayloadError(err error) {
	fmt.Println("Post: I am having problem generating the payload.", err.Error())
	fmt.Println("")
}

func doDelete(client *http.Client, p string, md []byte) {
	req, err := http.NewRequest(http.MethodDelete, p, bytes.NewBuffer(md))
	if err != nil {
		fmt.Println("Post:Delete: Problem connecting to Aegis Safe API endpoint URL.", err.Error())
		fmt.Println("")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := client.Do(req)
	if err != nil {
		fmt.Println("Post:Delete: Problem connecting to Aegis Safe API endpoint URL.", err.Error())
		fmt.Println("")
		return
	}
	respond(r)
}

func doPost(client *http.Client, p string, md []byte) {
	r, err := client.Post(p, "application/json", bytes.NewBuffer(md))
	if err != nil {
		fmt.Println("Post: Problem connecting to Aegis Safe API endpoint URL.", err.Error())
		fmt.Println("")
		return
	}
	respond(r)
}

func Post(workloadId, secret, namespace, backingStore string, useKubernetes bool,
	template string, format string, encrypt, deleteSecret, appendSecret bool,
) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	source, proceed := acquireSource(ctx)
	defer func() {
		if source == nil {
			return
		}
		err := source.Close()
		if err != nil {
			log.Println("Problem closing the workload source.")
		}
	}()
	if !proceed {
		return
	}

	authorizer := createAuthorizer()

	p, err := url.JoinPath(env.SafeEndpointUrl(), "/sentinel/v1/secrets")
	if err != nil {
		printEndpointError(err)
		return
	}

	tlsConfig := tlsconfig.MTLSClientConfig(source, source, authorizer)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	sr := newSecretUpsertRequest(workloadId, secret, namespace, backingStore,
		useKubernetes, template, format, encrypt, appendSecret)
	md, err := json.Marshal(sr)
	if err != nil {
		printPayloadError(err)
		return
	}

	if deleteSecret {
		doDelete(client, p, md)
		return
	}

	doPost(client, p, md)
}
