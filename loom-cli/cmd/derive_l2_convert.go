package cmd

import "github.com/ikadar/loom-cli/internal/formatter"

// convertContractsToFormatter converts local InterfaceContract slice to formatter types
func convertContractsToFormatter(contracts []InterfaceContract) []formatter.InterfaceContract {
	result := make([]formatter.InterfaceContract, len(contracts))
	for i, ic := range contracts {
		ops := make([]formatter.ContractOperation, len(ic.Operations))
		for j, op := range ic.Operations {
			inputSchema := make(map[string]formatter.SchemaField)
			for k, v := range op.InputSchema {
				inputSchema[k] = formatter.SchemaField{Type: v.Type, Required: v.Required}
			}
			outputSchema := make(map[string]formatter.SchemaField)
			for k, v := range op.OutputSchema {
				outputSchema[k] = formatter.SchemaField{Type: v.Type, Required: v.Required}
			}
			errors := make([]formatter.ContractError, len(op.Errors))
			for k, e := range op.Errors {
				errors[k] = formatter.ContractError{Code: e.Code, HTTPStatus: e.HTTPStatus, Message: e.Message}
			}
			ops[j] = formatter.ContractOperation{
				ID:             op.ID,
				Name:           op.Name,
				Method:         op.Method,
				Path:           op.Path,
				Description:    op.Description,
				InputSchema:    inputSchema,
				OutputSchema:   outputSchema,
				Errors:         errors,
				Preconditions:  op.Preconditions,
				Postconditions: op.Postconditions,
				RelatedACs:     op.RelatedACs,
				RelatedBRs:     op.RelatedBRs,
			}
		}
		events := make([]formatter.ContractEvent, len(ic.Events))
		for j, ev := range ic.Events {
			events[j] = formatter.ContractEvent{Name: ev.Name, Description: ev.Description, Payload: ev.Payload}
		}
		result[i] = formatter.InterfaceContract{
			ID:          ic.ID,
			ServiceName: ic.ServiceName,
			Purpose:     ic.Purpose,
			BaseURL:     ic.BaseURL,
			Operations:  ops,
			Events:      events,
			SecurityRequirements: formatter.SecurityRequirements{
				Authentication: ic.SecurityRequirements.Authentication,
				Authorization:  ic.SecurityRequirements.Authorization,
			},
		}
	}
	return result
}

// convertSharedTypesToFormatter converts local SharedType slice to formatter types
func convertSharedTypesToFormatter(types []SharedType) []formatter.SharedType {
	result := make([]formatter.SharedType, len(types))
	for i, st := range types {
		fields := make([]formatter.TypeField, len(st.Fields))
		for j, f := range st.Fields {
			fields[j] = formatter.TypeField{Name: f.Name, Type: f.Type, Constraints: f.Constraints}
		}
		result[i] = formatter.SharedType{Name: st.Name, Fields: fields}
	}
	return result
}

// convertAggregatesToFormatter converts local AggregateDesign slice to formatter types
func convertAggregatesToFormatter(aggregates []AggregateDesign) []formatter.AggregateDesign {
	result := make([]formatter.AggregateDesign, len(aggregates))
	for i, agg := range aggregates {
		invariants := make([]formatter.AggInvariant, len(agg.Invariants))
		for j, inv := range agg.Invariants {
			invariants[j] = formatter.AggInvariant{ID: inv.ID, Rule: inv.Rule, Enforcement: inv.Enforcement}
		}
		rootAttrs := make([]formatter.AggAttribute, len(agg.Root.Attributes))
		for j, attr := range agg.Root.Attributes {
			rootAttrs[j] = formatter.AggAttribute{Name: attr.Name, Type: attr.Type, Mutable: attr.Mutable}
		}
		entities := make([]formatter.AggEntity, len(agg.Entities))
		for j, ent := range agg.Entities {
			attrs := make([]formatter.AggAttribute, len(ent.Attributes))
			for k, attr := range ent.Attributes {
				attrs[k] = formatter.AggAttribute{Name: attr.Name, Type: attr.Type, Mutable: attr.Mutable}
			}
			entities[j] = formatter.AggEntity{Name: ent.Name, Identity: ent.Identity, Purpose: ent.Purpose, Attributes: attrs}
		}
		behaviors := make([]formatter.AggBehavior, len(agg.Behaviors))
		for j, b := range agg.Behaviors {
			behaviors[j] = formatter.AggBehavior{
				Name: b.Name, Command: b.Command, Preconditions: b.Preconditions,
				Postconditions: b.Postconditions, Emits: b.Emits,
			}
		}
		events := make([]formatter.AggEvent, len(agg.Events))
		for j, ev := range agg.Events {
			events[j] = formatter.AggEvent{Name: ev.Name, Payload: ev.Payload}
		}
		repoMethods := make([]formatter.RepoMethod, len(agg.Repository.Methods))
		for j, m := range agg.Repository.Methods {
			repoMethods[j] = formatter.RepoMethod{Name: m.Name, Params: m.Params, Returns: m.Returns}
		}
		extRefs := make([]formatter.AggExternalRef, len(agg.ExternalReferences))
		for j, ref := range agg.ExternalReferences {
			extRefs[j] = formatter.AggExternalRef{Aggregate: ref.Aggregate, Via: ref.Via, Type: ref.Type}
		}
		result[i] = formatter.AggregateDesign{
			ID:           agg.ID,
			Name:         agg.Name,
			Purpose:      agg.Purpose,
			Invariants:   invariants,
			Root:         formatter.AggRoot{Entity: agg.Root.Entity, Identity: agg.Root.Identity, Attributes: rootAttrs},
			Entities:     entities,
			ValueObjects: agg.ValueObjects,
			Behaviors:    behaviors,
			Events:       events,
			Repository: formatter.AggRepository{
				Name: agg.Repository.Name, Methods: repoMethods,
				LoadStrategy: agg.Repository.LoadStrategy, Concurrency: agg.Repository.Concurrency,
			},
			ExternalReferences: extRefs,
		}
	}
	return result
}

// convertSequencesToFormatter converts local SequenceDesign slice to formatter types
func convertSequencesToFormatter(sequences []SequenceDesign) []formatter.SequenceDesign {
	result := make([]formatter.SequenceDesign, len(sequences))
	for i, seq := range sequences {
		participants := make([]formatter.SeqParticipant, len(seq.Participants))
		for j, p := range seq.Participants {
			participants[j] = formatter.SeqParticipant{Name: p.Name, Type: p.Type}
		}
		steps := make([]formatter.SequenceStep, len(seq.Steps))
		for j, s := range seq.Steps {
			steps[j] = formatter.SequenceStep{
				Step: s.Step, Actor: s.Actor, Action: s.Action, Target: s.Target,
				Data: s.Data, Returns: s.Returns, Event: s.Event,
			}
		}
		exceptions := make([]formatter.SequenceException, len(seq.Exceptions))
		for j, ex := range seq.Exceptions {
			exceptions[j] = formatter.SequenceException{Condition: ex.Condition, Step: ex.Step, Handling: ex.Handling}
		}
		result[i] = formatter.SequenceDesign{
			ID:           seq.ID,
			Name:         seq.Name,
			Description:  seq.Description,
			Trigger:      formatter.SequenceTrigger{Type: seq.Trigger.Type, Description: seq.Trigger.Description},
			Participants: participants,
			Steps:        steps,
			Outcome: formatter.SequenceOutcome{
				Success: seq.Outcome.Success, StateChanges: seq.Outcome.StateChanges,
			},
			Exceptions: exceptions,
			RelatedACs: seq.RelatedACs,
			RelatedBRs: seq.RelatedBRs,
		}
	}
	return result
}

// convertTablesToFormatter converts local DataTable slice to formatter types
func convertTablesToFormatter(tables []DataTable) []formatter.DataTable {
	result := make([]formatter.DataTable, len(tables))
	for i, tbl := range tables {
		fields := make([]formatter.DataField, len(tbl.Fields))
		for j, f := range tbl.Fields {
			fields[j] = formatter.DataField{Name: f.Name, Type: f.Type, Constraints: f.Constraints, Default: f.Default}
		}
		indexes := make([]formatter.DataIndex, len(tbl.Indexes))
		for j, idx := range tbl.Indexes {
			indexes[j] = formatter.DataIndex{Name: idx.Name, Columns: idx.Columns}
		}
		fks := make([]formatter.DataForeignKey, len(tbl.ForeignKeys))
		for j, fk := range tbl.ForeignKeys {
			fks[j] = formatter.DataForeignKey{Columns: fk.Columns, References: fk.References, OnDelete: fk.OnDelete}
		}
		constraints := make([]formatter.DataConstraint, len(tbl.CheckConstraints))
		for j, c := range tbl.CheckConstraints {
			constraints[j] = formatter.DataConstraint{Name: c.Name, Expression: c.Expression}
		}
		result[i] = formatter.DataTable{
			ID:               tbl.ID,
			Name:             tbl.Name,
			Aggregate:        tbl.Aggregate,
			Purpose:          tbl.Purpose,
			Fields:           fields,
			PrimaryKey:       formatter.DataPrimaryKey{Columns: tbl.PrimaryKey.Columns},
			Indexes:          indexes,
			ForeignKeys:      fks,
			CheckConstraints: constraints,
		}
	}
	return result
}

// convertEnumsToFormatter converts local DataEnum slice to formatter types
func convertEnumsToFormatter(enums []DataEnum) []formatter.DataEnum {
	result := make([]formatter.DataEnum, len(enums))
	for i, e := range enums {
		result[i] = formatter.DataEnum{Name: e.Name, Values: e.Values}
	}
	return result
}
