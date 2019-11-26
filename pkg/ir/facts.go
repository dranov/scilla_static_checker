package ir

import (
	//"errors"
	"fmt"
	"os"
	"path"
	//"gitlab.chainsecurity.com/ChainSecurity/common/scilla_static/pkg/ast"
	//"strings"
)

type FactsDumper struct {
	visited         map[Node]bool
	idToPrefixID    map[uint64]string
	procFacts       []string
	dataVarFacts    []string
	transitionFacts []string
	procedureFacts  []string
	unitFacts       []string
	planFacts       []string
	sendFacts       []string
	acceptFacts     []string
	saveFacts       []string
	loadFacts       []string
	appDDFacts      []string
	appTDFacts      []string
	appTTFacts      []string
	argumentFacts   []string
	absDDFacts      []string
	msgFacts        []string
	msgTypeFacts    []string
	msgDataFacts    []string
	strFacts        []string
	strTypeFacts    []string
	pickDataFacts   []string
	dataCaseFacts   []string
	jumpFacts       []string
	procCaseFacts   []string
	callProcFacts   []string
	pickProcFacts   []string
	natFacts        []string
	natTypeFacts    []string
	intFacts        []string
	intTypeFacts    []string
	rawFacts        []string
	rawTypeFacts    []string
	bindFacts       []string
	condFacts       []string
	condBindFacts   []string
	typeVarFacts    []string
	enumFacts       []string
	enumTypeFacts   []string
	Facts           []string
}

func (fd *FactsDumper) Visit(node Node, prev Node) Visitor {

	visited, ok := fd.visited[node]
	if ok && visited {
		return nil
	}

	// DEBUG related check
	if node.ID() == 0 {
		panic(fmt.Sprintf("ID was not assigned %T", node))
	}

	fd.unitFacts = append(fd.unitFacts, fmt.Sprintf("%d", node.ID()))

	switch n := node.(type) {
	case *Proc:
		//prefixID := fmt.Sprintf("Proc%d", n.ID())
		//fd.idToPrefixID[n.ID()] = prefixID
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.ProcName)
		fd.procFacts = append(fd.procFacts, fact)

		for i, u := range n.Plan {
			fact := fmt.Sprintf("plan_%d_%d\t%d\t%d", n.ID(), i, n.ID(), u.ID())
			fd.planFacts = append(fd.planFacts, fact)
		}

		if n.Jump != nil {
			fact = fmt.Sprintf("%d\t%d", n.ID(), n.Jump.ID())
			fd.jumpFacts = append(fd.jumpFacts, fact)
		}

	case *DataVar:
		fact := fmt.Sprintf("%d", n.ID())
		fd.dataVarFacts = append(fd.dataVarFacts, fact)
	case *Send:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Data.ID())
		fd.sendFacts = append(fd.sendFacts, fact)
	case *Accept:
		fact := fmt.Sprintf("%d", n.ID())
		fd.acceptFacts = append(fd.acceptFacts, fact)
	case *Save:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.saveFacts = append(fd.saveFacts, fact)
	case *Load:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Slot)
		fd.loadFacts = append(fd.loadFacts, fact)
	case *AppDD:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appDDFacts = append(fd.appDDFacts, fact)
	case *AppTD:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appTDFacts = append(fd.appTDFacts, fact)
	case *AppTT:
		for i, u := range n.Args {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.To.ID())
		fd.appTTFacts = append(fd.appTTFacts, fact)
	case *AbsDD:
		for i, u := range n.Vars {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Term.ID())
		fd.absDDFacts = append(fd.absDDFacts, fact)
	case *Msg:
		for k, v := range n.Data {
			fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), v.ID(), k)
			fd.msgDataFacts = append(fd.msgDataFacts, fact)
		}
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.MsgType.ID())
		fd.msgFacts = append(fd.msgFacts, fact)
	case *MsgType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.msgTypeFacts = append(fd.msgTypeFacts, fact)
	case *Str:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.StrType, n.Data)
		fd.strFacts = append(fd.strFacts, fact)
	case *StrType:
		fact := fmt.Sprintf("%d", n.ID())
		fd.strTypeFacts = append(fd.strTypeFacts, fact)
	case *PickData:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.From)
		fd.pickDataFacts = append(fd.pickDataFacts, fact)
		for i, u := range n.With {
			argFact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.argumentFacts = append(fd.argumentFacts, argFact)
		}

	//case *PickProc:
	case *DataCase:

		body, ok := n.Body.(Node)
		if !ok {
			panic(fmt.Sprintf("Node is IDNode %T", body))
		}

		fact := fmt.Sprintf("%d\t%d\t%d\t%d", n.ID(), prev.ID(), n.Bind.ID(), body.ID())
		fd.dataCaseFacts = append(fd.dataCaseFacts, fact)
	case *Bind:
		var condID int64
		condID = -1
		if n.Cond != nil {
			condID = n.Cond.ID()
		}
		fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), n.BindType.ID(), condID)
		fd.bindFacts = append(fd.bindFacts, fact)
	case *Cond:
		fact := fmt.Sprintf("%d\t%s", n.ID(), n.Case)
		fd.condFacts = append(fd.condFacts, fact)

		for i, u := range n.Data {
			fact := fmt.Sprintf("%d\t%d\t%d", n.ID(), u.ID(), i)
			fd.condBindFacts = append(fd.condBindFacts, fact)
		}
	case *Nat:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.NatType.ID(), n.Data)
		fd.natFacts = append(fd.natFacts, fact)
	case *NatType:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Size)
		fd.natTypeFacts = append(fd.natTypeFacts, fact)

	case *Int:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.IntType.ID(), n.Data)
		fd.intFacts = append(fd.intFacts, fact)
	case *IntType:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Size)
		fd.intTypeFacts = append(fd.intTypeFacts, fact)

	case *Raw:
		fact := fmt.Sprintf("%d\t%d\t%s", n.ID(), n.RawType.ID(), n.Data)
		fd.rawFacts = append(fd.rawFacts, fact)
	case *RawType:
		fact := fmt.Sprintf("%d\t%d", n.ID(), n.Size)
		fd.rawTypeFacts = append(fd.rawTypeFacts, fact)

	default:
		fmt.Printf("+ %T %d\n", node, node.ID())
	}

	fd.visited[node] = true

	return fd
}

func DumpFacts(builder *CFGBuilder) {
	fd := FactsDumper{
		visited:       map[Node]bool{},
		idToPrefixID:  map[uint64]string{},
		procFacts:     []string{},
		dataVarFacts:  []string{},
		unitFacts:     []string{},
		planFacts:     []string{},
		sendFacts:     []string{},
		acceptFacts:   []string{},
		saveFacts:     []string{},
		loadFacts:     []string{},
		appDDFacts:    []string{},
		appTDFacts:    []string{},
		appTTFacts:    []string{},
		argumentFacts: []string{},
		absDDFacts:    []string{},
		msgFacts:      []string{},
		msgDataFacts:  []string{},
		strFacts:      []string{},
		strTypeFacts:  []string{},
		callProcFacts: []string{},
		pickProcFacts: []string{},
		pickDataFacts: []string{},
		dataCaseFacts: []string{},
		procCaseFacts: []string{},
		natFacts:      []string{},
		natTypeFacts:  []string{},
		intFacts:      []string{},
		intTypeFacts:  []string{},
		rawFacts:      []string{},
		rawTypeFacts:  []string{},
		bindFacts:     []string{},
		condFacts:     []string{},
		condBindFacts: []string{},
		typeVarFacts:  []string{},
		enumFacts:     []string{},
		enumTypeFacts: []string{},
		Facts:         []string{},
	}
	for tName, t := range builder.Transitions {
		fmt.Println("Transition", tName)
		fact := fmt.Sprintf("%d", t.ID())
		fd.transitionFacts = append(fd.transitionFacts, fact)
		Walk(&fd, t, nil)
	}
	for pName, p := range builder.Procedures {
		fmt.Println("Procedure", pName)
		fact := fmt.Sprintf("%d", p.ID())
		fd.procedureFacts = append(fd.procedureFacts, fact)
		Walk(&fd, p, nil)
	}

	fileToFacts := map[string][]string{
		"proc":       fd.procFacts,
		"dataVar":    fd.dataVarFacts,
		"transition": fd.transitionFacts,
		"procedure":  fd.procedureFacts,
		"unit":       fd.unitFacts,
		"plan":       fd.planFacts,
		"send":       fd.sendFacts,
		"accept":     fd.acceptFacts,
		"save":       fd.saveFacts,
		"load":       fd.loadFacts,
		"appDD":      fd.appDDFacts,
		"appTD":      fd.appTDFacts,
		"appTT":      fd.appTTFacts,
		"argument":   fd.argumentFacts,
		"absDD":      fd.absDDFacts,
		"msg":        fd.msgFacts,
		"msgData":    fd.msgDataFacts,
		"str":        fd.strFacts,
		"strType":    fd.strTypeFacts,
		"jump":       fd.jumpFacts,
		"callProc":   fd.callProcFacts,
		"pickProc":   fd.pickProcFacts,
		"pickData":   fd.pickDataFacts,
		"dataCase":   fd.dataCaseFacts,
		"nat":        fd.natFacts,
		"natType":    fd.natTypeFacts,
		"int":        fd.intFacts,
		"intType":    fd.intTypeFacts,
		"raw":        fd.rawFacts,
		"rawType":    fd.rawTypeFacts,
		"bind":       fd.bindFacts,
		"cond":       fd.condFacts,
		"condBind":   fd.condBindFacts,
		"typeVar":    fd.typeVarFacts,
		"enum":       fd.enumFacts,
		"enumType":   fd.enumTypeFacts,
	}

	analysisFolder := "./souffle_analysis"
	if _, err := os.Stat(analysisFolder); os.IsNotExist(err) {
		err = os.Mkdir(analysisFolder, 0700)
		if err != nil {
			panic(err)
		}
	}

	factsInFolder := path.Join(analysisFolder, "facts_in")

	if _, err := os.Stat(factsInFolder); !os.IsNotExist(err) {
		err = os.RemoveAll(factsInFolder)
		if err != nil {
			panic(err)
		}
	}
	err := os.Mkdir(factsInFolder, 0700)
	if err != nil {
		panic(err)
	}

	for fileName, lines := range fileToFacts {
		filePath := path.Join(factsInFolder, fmt.Sprintf("%s.facts", fileName))
		f, err := os.Create(filePath)
		if err != nil {
			f.Close()
			panic(err)
		}
		for _, line := range lines {
			fmt.Fprintln(f, line)
			if err != nil {
				panic(err)
			}
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
