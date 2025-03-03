package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.34

import (
	"context"
	"errors"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/graph/generated"
	"github.com/BuxOrg/bux/utils"
	"github.com/mrz1836/go-datastore"
)

// Xpub is the resolver for the xpub field.
func (r *mutationResolver) Xpub(ctx context.Context, xpub string, metadata bux.Metadata) (*bux.Xpub, error) {
	// including admin check
	c, err := GetConfigFromContextAdmin(ctx)
	if err != nil {
		return nil, err
	}

	var existingXpub *bux.Xpub
	existingXpub, err = c.Services.Bux.GetXpub(ctx, xpub)
	if err != nil && !errors.Is(err, bux.ErrMissingXpub) {
		return nil, err
	}
	if existingXpub != nil {
		return nil, errors.New("xpub already exists")
	}

	opts := c.Services.Bux.DefaultModelOptions()
	for key, value := range metadata {
		opts = append(opts, bux.WithMetadata(key, value))
	}

	// Create a new xPub
	var xPub *bux.Xpub
	if xPub, err = c.Services.Bux.NewXpub(
		ctx, xpub, opts...,
	); err != nil {
		return nil, err
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

// XpubMetadata is the resolver for the xpub_metadata field.
func (r *mutationResolver) XpubMetadata(ctx context.Context, metadata bux.Metadata) (*bux.Xpub, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var xPub *bux.Xpub
	xPub, err = c.Services.Bux.UpdateXpubMetadata(ctx, c.XPubID, metadata)
	if err != nil {
		return nil, err
	}

	if !c.Signed || c.XPub == "" {
		xPub.RemovePrivateData()
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

// AccessKey is the resolver for the access_key field.
func (r *mutationResolver) AccessKey(ctx context.Context, metadata bux.Metadata) (*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	// Create a new accessKey
	var accessKey *bux.AccessKey
	if accessKey, err = c.Services.Bux.NewAccessKey(
		ctx,
		c.XPub,
		bux.WithMetadatas(metadata),
	); err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKey).(*bux.AccessKey), nil
}

// AccessKeyRevoke is the resolver for the access_key_revoke field.
func (r *mutationResolver) AccessKeyRevoke(ctx context.Context, id *string) (*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	// Revoke an accessKey
	var accessKey *bux.AccessKey
	if accessKey, err = c.Services.Bux.RevokeAccessKey(
		ctx,
		c.XPub,
		*id,
	); err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKey).(*bux.AccessKey), nil
}

// Transaction is the resolver for the transaction field.
func (r *mutationResolver) Transaction(ctx context.Context, hex string, draftID *string, metadata bux.Metadata) (*bux.Transaction, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()
	for key, value := range metadata {
		opts = append(opts, bux.WithMetadata(key, value))
	}

	ref := ""
	if draftID != nil {
		ref = *draftID
	}

	var transaction *bux.Transaction
	transaction, err = c.Services.Bux.RecordTransaction(
		ctx, c.XPub, hex, ref, opts...,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrDuplicateKey) {
			var txID string
			txID, err = utils.GetTransactionIDFromHex(hex)
			if err != nil {
				return nil, err
			}

			transaction, err = c.Services.Bux.GetTransaction(ctx, c.XPub, txID)
			if err != nil {
				return nil, err
			}

			// record the metadata is being added to the transaction
			if len(metadata) > 0 {
				xPubID := utils.Hash(c.XPub)
				if transaction.XpubMetadata == nil {
					transaction.XpubMetadata = make(bux.XpubMetadata)
				}
				if transaction.XpubMetadata[xPubID] == nil {
					transaction.XpubMetadata[xPubID] = make(bux.Metadata)
				}
				for key, value := range metadata {
					transaction.XpubMetadata[xPubID][key] = value
				}
				err = transaction.Save(ctx)
				if err != nil {
					return nil, err
				}
				// set metadata to the xpub metadata - is removed after Save
				transaction.Metadata = transaction.XpubMetadata[xPubID]
			}

			return transaction, nil
		}
		return nil, err
	}

	return bux.DisplayModels(transaction).(*bux.Transaction), nil
}

// TransactionMetadata is the resolver for the transaction_metadata field.
func (r *mutationResolver) TransactionMetadata(ctx context.Context, id string, metadata bux.Metadata) (*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx *bux.Transaction
	tx, err = c.Services.Bux.UpdateTransactionMetadata(ctx, c.XPubID, id, metadata)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return bux.DisplayModels(tx).(*bux.Transaction), nil
}

// NewTransaction is the resolver for the new_transaction field.
func (r *mutationResolver) NewTransaction(ctx context.Context, transactionConfig bux.TransactionConfig, metadata bux.Metadata) (*bux.DraftTransaction, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	opts := c.Services.Bux.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	var draftTransaction *bux.DraftTransaction
	draftTransaction, err = c.Services.Bux.NewTransaction(ctx, c.XPub, &transactionConfig, opts...)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(draftTransaction).(*bux.DraftTransaction), nil
}

// Destination is the resolver for the destination field.
func (r *mutationResolver) Destination(ctx context.Context, destinationType *string, metadata bux.Metadata) (*bux.Destination, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var useDestinationType string
	if destinationType != nil {
		useDestinationType = *destinationType
	} else {
		useDestinationType = utils.ScriptTypePubKeyHash
	}

	opts := c.Services.Bux.DefaultModelOptions()
	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	var destination *bux.Destination
	destination, err = c.Services.Bux.NewDestination(
		ctx,
		c.XPub,
		utils.ChainExternal,
		useDestinationType,
		true, // monitor this address as it was created by request of a user to share
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

// DestinationMetadata is the resolver for the destination_metadata field.
func (r *mutationResolver) DestinationMetadata(ctx context.Context, id *string, address *string, lockingScript *string, metadata bux.Metadata) (*bux.Destination, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var destination *bux.Destination
	if id != nil {
		destination, err = c.Services.Bux.UpdateDestinationMetadataByID(
			ctx,
			c.XPubID,
			*id,
			metadata,
		)
	} else if address != nil {
		destination, err = c.Services.Bux.UpdateDestinationMetadataByAddress(
			ctx,
			c.XPubID,
			*address,
			metadata,
		)
	} else if lockingScript != nil {
		destination, err = c.Services.Bux.UpdateDestinationMetadataByLockingScript(
			ctx,
			c.XPubID,
			*lockingScript,
			metadata,
		)
	}
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

// UtxosUnreserve is the resolver for the utxos_unreserve field.
func (r *mutationResolver) UtxosUnreserve(ctx context.Context, draftID string) (*bool, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Services.Bux.UnReserveUtxos(
		ctx,
		c.XPubID,
		draftID,
	)

	var success bool
	success = err != nil
	return &success, err
}

// Xpub is the resolver for the xpub field.
func (r *queryResolver) Xpub(ctx context.Context) (*bux.Xpub, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var xPub *bux.Xpub
	xPub, err = c.Services.Bux.GetXpubByID(ctx, c.XPubID)
	if err != nil {
		return nil, err
	}

	if !c.Signed || c.XPub == "" {
		xPub.RemovePrivateData()
	}

	return bux.DisplayModels(xPub).(*bux.Xpub), nil
}

// AccessKey is the resolver for the access_key field.
func (r *queryResolver) AccessKey(ctx context.Context, key string) (*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var accessKey *bux.AccessKey
	accessKey, err = c.Services.Bux.GetAccessKey(ctx, c.XPubID, key)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKey).(*bux.AccessKey), nil
}

// AccessKeys is the resolver for the access_keys field.
func (r *queryResolver) AccessKeys(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.AccessKey, error) {
	c, err := GetConfigFromContextSigned(ctx)
	if err != nil {
		return nil, err
	}

	var accessKeys []*bux.AccessKey
	accessKeys, err = c.Services.Bux.GetAccessKeysByXPubID(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions), params)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(accessKeys).([]*bux.AccessKey), nil
}

// AccessKeysCount is the resolver for the access_keys_count field.
func (r *queryResolver) AccessKeysCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetAccessKeysByXPubIDCount(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions))
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// Transaction is the resolver for the transaction field.
func (r *queryResolver) Transaction(ctx context.Context, id string) (*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx *bux.Transaction
	tx, err = c.Services.Bux.GetTransaction(ctx, c.XPubID, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return bux.DisplayModels(tx).(*bux.Transaction), nil
}

// Transactions is the resolver for the transactions field.
func (r *queryResolver) Transactions(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Transaction, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var tx []*bux.Transaction
	tx, err = c.Services.Bux.GetTransactionsByXpubID(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions), params)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(tx).([]*bux.Transaction), nil
}

// TransactionsCount is the resolver for the transactions_count field.
func (r *queryResolver) TransactionsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetTransactionsByXpubIDCount(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions))
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// Destination is the resolver for the destination field.
func (r *queryResolver) Destination(ctx context.Context, id *string, address *string, lockingScript *string) (*bux.Destination, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var destination *bux.Destination
	if id != nil {
		destination, err = c.Services.Bux.GetDestinationByID(ctx, c.XPubID, *id)
	} else if address != nil {
		destination, err = c.Services.Bux.GetDestinationByAddress(ctx, c.XPubID, *address)
	} else if lockingScript != nil {
		destination, err = c.Services.Bux.GetDestinationByLockingScript(ctx, c.XPubID, *lockingScript)
	} else {
		return nil, bux.ErrMissingFieldID
	}
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destination).(*bux.Destination), nil
}

// Destinations is the resolver for the destinations field.
func (r *queryResolver) Destinations(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Destination, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var destinations []*bux.Destination
	destinations, err = c.Services.Bux.GetDestinationsByXpubID(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions), params)
	if err != nil {
		return nil, err
	}

	return bux.DisplayModels(destinations).([]*bux.Destination), nil
}

// DestinationsCount is the resolver for the destinations_count field.
func (r *queryResolver) DestinationsCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var count int64
	count, err = c.Services.Bux.GetDestinationsByXpubIDCount(ctx, c.XPubID, &metadata, ConditionsParseGraphQL(conditions))
	if err != nil {
		return nil, err
	}

	return &count, nil
}

// Utxo is the resolver for the utxo field.
func (r *queryResolver) Utxo(ctx context.Context, txID string, outputIndex uint32) (*bux.Utxo, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var utxo *bux.Utxo
	if utxo, err = c.Services.Bux.GetUtxo(
		ctx,
		c.XPubID,
		txID,
		outputIndex,
	); err != nil {
		return nil, err
	}

	return utxo, nil
}

// Utxos is the resolver for the utxos field.
func (r *queryResolver) Utxos(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}, params *datastore.QueryParams) ([]*bux.Utxo, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var utxos []*bux.Utxo
	if utxos, err = c.Services.Bux.GetUtxosByXpubID(
		ctx,
		c.XPubID,
		&metadata,
		&conditions,
		params,
	); err != nil {
		return nil, err
	}

	return utxos, nil
}

// UtxosCount is the resolver for the utxos_count field.
func (r *queryResolver) UtxosCount(ctx context.Context, metadata bux.Metadata, conditions map[string]interface{}) (*int64, error) {
	c, err := GetConfigFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = conditions
	}
	// force the xpub_id of the current user on query
	dbConditions["xpub_id"] = c.XPubID

	var count int64
	if count, err = c.Services.Bux.GetUtxosCount(
		ctx,
		&metadata,
		&dbConditions,
	); err != nil {
		return nil, err
	}

	return &count, nil
}

// Inputs is the resolver for the inputs field.
func (r *transactionConfigInputResolver) Inputs(ctx context.Context, obj *bux.TransactionConfig, data []map[string]interface{}) error {
	// do nothing with inputs
	return nil
}

// ExpiresIn is the resolver for the expires_in field.
func (r *transactionConfigInputResolver) ExpiresIn(ctx context.Context, obj *bux.TransactionConfig, data *uint64) error {
	obj.ExpiresIn = time.Duration(*data) * time.Second
	return nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// TransactionConfigInput returns generated.TransactionConfigInputResolver implementation.
func (r *Resolver) TransactionConfigInput() generated.TransactionConfigInputResolver {
	return &transactionConfigInputResolver{r}
}

type (
	mutationResolver               struct{ *Resolver }
	queryResolver                  struct{ *Resolver }
	transactionConfigInputResolver struct{ *Resolver }
)
