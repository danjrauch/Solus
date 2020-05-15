package main

type PropertyGraph interface {
	// CRUD Operations
	// AddNode()
	// AddEdge()
	// DeleteNode()
	// DeleteEdge()
	// GetNode(id)
	// GetEdge(id)

	// Operations over Schema
	// GetNodeLabels()
	// GetEdgeLabels()
	// ScanNodes(type)
	// ScanEdges(type)
	// GetNodeLabel(nodeId)
	// GetEdgeLabel(edgeId)

	// Operations over Data
	// GetNodeAttribute(nodeId, attribute)
	// GetEdgeAttribute(edgeId, attribute)
	// SelectNodes(type, attribute, val)
	// SelectEdges(type, attribute, val)

	// Operations over Relationships
	// GetNeighbors(type, id)
	// GetRelated(type, id)
}

type id int

type propertyGraph struct {
	label2Id         map[string]id
	id2Label         map[id]string
	attribute2Id     map[string]id
	id2Attribute     map[id]string
	nodeSchema       map[id][]id
	edgeSchema       map[id][]id
	attributeValues  map[id][]interface{}
	attributeCounter id
	labelCounter     id
	nodeCounter      id
	edgeCounter      id
	nodes            map[id]node
	edges            map[id]edge
	strict           bool
}

type node struct {
	id         id
	label      id
	attributes map[id][]interface{}
}

type edge struct {
	id         id
	from       id
	to         id
	label      id
	attributes map[id][]interface{}
}
