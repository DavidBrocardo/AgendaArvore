//Autores:
//- Leonardo Balan
//- David Brocardo

// Agenda de Contatos indexada por Árvore-B
//AgendaArvoreBQ1.go

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Constantes ultilizadas
const t = 3                  // Arvore terá um grau 3
const ArqDados = "Dados.txt" //  Constante para o arquivo que seram armazenados os dados
const ArqInd = "Indices.txt" //  Constante para o arquivo que seram armazenados a chave e posição

// Constantes para garantir que não execeda o tamanho proposto
const (
	MaxNome     = 30
	MaxEndereco = 50
	MaxTelefone = 15
)

type DataType struct {
	nome   string
	indice int64
}

// Struct para Nó
type BTreeNode struct {
	leaf     bool
	keys     []DataType
	children []*BTreeNode
}

// Para verificar se é folha ou não
func InitNode(leaf bool) *BTreeNode {
	return &BTreeNode{
		leaf:     leaf,
		keys:     []DataType{},
		children: []*BTreeNode{},
	}
}

// struct da arvore B
type BTree struct {
	root *BTreeNode
}

// inicialização
func Init() *BTree {
	return &BTree{
		root: InitNode(true),
	}
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função para Printar Elementos
func (node *BTreeNode) Listar_Ord(dados_arq string) {
	file, _ := os.Open(dados_arq)
	defer file.Close()

	if !node.leaf { // Se não for um nó folha, ou seja, é um nó interno da árvore
		for i := 0; i < len(node.keys); i++ {
			node.children[i].Listar_Ord(dados_arq) // Chama recursivamente Listar_Ord para os filhos à esquerda do nó atual
			node.keys[i].Printar_Contato(file)     // Imprime o contato do índice i no arquivo de dados
		}
		node.children[len(node.keys)].Listar_Ord(dados_arq) // Chama recursivamente Listar_Ord para o último filho à direita do nó atual
	} else {
		for j := 0; j < len(node.keys); j++ { // Se for um nó folha, ou seja, uma folha da árvore
			node.keys[j].Printar_Contato(file) // Imprime o contato do índice j no arquivo de dados
		}
	}
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função para Remover um Contato
func (tree *BTree) Remove_Contato(nome string, file_data_nome string) {
	file, _ := os.OpenFile(file_data_nome, os.O_RDWR, 0644)

	no := tree.ProcuraNome(nome) // Encontra o nó do contato a ser removido

	if no == nil {
		fmt.Println(" ")
		fmt.Println("Nome Inválido.")
		fmt.Println("Não foi possível realizar a remoção.")
		fmt.Println(" ")
		return
	}

	file.Seek(no.indice, 0)

	buffer := make([]byte, 95)
	file.Read(buffer)

	var marca_l int
	marca_l = 0

	for i := 0; i < 95; i++ {
		if marca_l == 2 && string(buffer[i]) == "\n" { // Encontra o byte depois de 2 \n (onde está a marca)
			marca_l = i + 1
			break
		}

		if string(buffer[i]) == "\n" {
			marca_l += 1
		}
	}

	new_byte := byte('0')

	file.WriteAt([]byte{new_byte}, int64(marca_l)+no.indice) // Altera a marca para 0

	file.Close()
	fmt.Println(" ")
	fmt.Println("Contato Removido ")
	fmt.Println(" ")
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função para Recuperar os Contatos da Lixeira
func RestauraLixeira(file_data_nome string) {
	file, _ := os.OpenFile(file_data_nome, os.O_RDWR, 0644)

	defer file.Close()

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, _ = file.Read(buffer)

	var marca_l, end_of_file, i int64
	marca_l = 2
	end_of_file = FindIndEOF(file) // Encontra o índice do fim do arquivo

	for i = 0; int64(i) < end_of_file; i++ {

		if int(buffer[i]) == 10 {
			marca_l += 1
		}

		// Realiza um salto para encontrar o byte que indica o contato como "apagado" (lixeira)
		if marca_l == 5 && int(buffer[i]) == 10 {
			new_byte := byte('1')
			file.WriteAt([]byte{new_byte}, int64(i+1)) // Altera o byte para '1' (contato apagado)
			marca_l = 0
		}

		if int64(i) == end_of_file {
			return
		}
	}
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função para Alterar o Nome
func (tree *BTree) Altera_nome(nome, novo_nome, dados_arq, ind_arq string) { // Altera o nome do contato
	file, err := os.OpenFile(dados_arq, os.O_RDWR, 0644)
	defer file.Close()

	var file_size int64
	file_size, _ = file.Seek(0, 2)
	no := tree.ProcuraNome(nome)

	if no == nil {
		fmt.Println("Nome Invalido.")
		return
	}

	file.Seek(no.indice, 0)
	buffer_find_ind := make([]byte, 95)
	file.Read(buffer_find_ind)

	var next_line int64

	for i := 0; i < 95; i++ {
		if string(buffer_find_ind[i]) == "\n" {
			next_line = no.indice + int64(i+1)
			break
		}
	}

	file.Seek(0, 0)
	buf_ant := make([]byte, no.indice)
	file.Read(buf_ant)

	novo_nome_b := []byte(novo_nome)
	buf_ant = append(buf_ant, novo_nome_b...)

	buf_post := make([]byte, file_size-next_line)

	var bytesRead int
	var totalBytesRead int

	for totalBytesRead < len(buf_post) {
		bytesRead, err = file.Read(buf_post[totalBytesRead:])
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}
		totalBytesRead += bytesRead
	}

	file.Seek(next_line, 0)
	file.Read(buf_post)

	aux_file, err := os.OpenFile("aux.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer aux_file.Close()

	aux_file.WriteString(string(buf_ant))
	aux_file.WriteString("\n")
	aux_file.WriteString(string(buf_post))

	Atualiza_indices("aux.txt", ind_arq)
	os.Rename("aux.txt", dados_arq)

	fmt.Println("Contato alterado.")

}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função para Alterar o Endereço
func (tree *BTree) Altera_end(nome, novo_end, dados_arq, ind_arq string) {
	file, err := os.OpenFile(dados_arq, os.O_RDWR, 0644)
	defer file.Close()

	var file_size int64
	file_size, _ = file.Seek(0, 2)
	no := tree.ProcuraNome(nome)

	if no == nil {
		fmt.Println("Nome Invalido.")
		return
	}

	file.Seek(no.indice, 0)
	buffer_find_ind := make([]byte, 95)
	file.Read(buffer_find_ind)

	var marca_l int
	var right_place, next_line int64
	marca_l = 0

	for i := 0; i < 95; i++ {
		if marca_l != 3 && string(buffer_find_ind[i]) == "\n" {
			marca_l += 1
		}

		if marca_l == 2 && string(buffer_find_ind[i]) == "\n" {
			right_place = no.indice + int64(i)
			marca_l += 1
			i += 1
		}

		if marca_l == 3 && string(buffer_find_ind[i]) == "\n" {
			next_line = no.indice + int64(i+1)
			break
		}
	}

	file.Seek(0, 0)
	buf_ant := make([]byte, right_place+1)
	file.Read(buf_ant)

	novo_end_b := []byte(novo_end)
	buf_ant = append(buf_ant, novo_end_b...)

	buf_post := make([]byte, file_size-next_line)

	var bytesRead int
	var totalBytesRead int

	for totalBytesRead < len(buf_post) {
		bytesRead, err = file.Read(buf_post[totalBytesRead:])
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}
		totalBytesRead += bytesRead
	}

	file.Seek(next_line, 0)
	file.Read(buf_post)

	aux_file, err := os.OpenFile("aux.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer aux_file.Close()

	aux_file.WriteString(string(buf_ant))
	aux_file.WriteString("\n")
	aux_file.WriteString(string(buf_post))

	Atualiza_indices("aux.txt", ind_arq)
	os.Rename("aux.txt", dados_arq)

	fmt.Println("Contato alterado.")

}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função para Alterar o Telefone
func (tree *BTree) Alterar_tel(nome, novo_tel, dados_arq, ind_arq string) {
	file, err := os.OpenFile(dados_arq, os.O_RDWR, 0644)
	defer file.Close()

	var file_size int64
	file_size, _ = file.Seek(0, 2)
	no := tree.ProcuraNome(nome)

	if no == nil {
		fmt.Println("Nome Invalido.")
		return
	}

	file.Seek(no.indice, 0)
	buffer_find_ind := make([]byte, 95) 
	file.Read(buffer_find_ind)         

	var marca_l int
	var right_place, next_line int64
	marca_l = 0

	for i := 0; i < 95; i++ {
		if marca_l != 2 && string(buffer_find_ind[i]) == "\n" {
			marca_l += 1
		}

		if marca_l == 1 && string(buffer_find_ind[i]) == "\n" {
			right_place = no.indice + int64(i) 
			marca_l += 1
			i += 1
		}

		if marca_l == 2 && string(buffer_find_ind[i]) == "\n" {
			next_line = no.indice + int64(i+1) 
			break
		}
	}

	file.Seek(0, 0)
	buf_ant := make([]byte, right_place+1) 
	file.Read(buf_ant)

	novo_tel_b := []byte(novo_tel)
	buf_ant = append(buf_ant, novo_tel_b...) 

	buf_post := make([]byte, file_size-next_line) 

	var bytesRead int
	var totalBytesRead int

	for totalBytesRead < len(buf_post) {
		bytesRead, err = file.Read(buf_post[totalBytesRead:])
		if err != nil {
			fmt.Println("Erro ao ler o arquivo:", err)
			return
		}
		totalBytesRead += bytesRead
	}

	file.Seek(next_line, 0)
	file.Read(buf_post)

	aux_file, err := os.OpenFile("aux.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer aux_file.Close()

	aux_file.WriteString(string(buf_ant))
	aux_file.WriteString("\n")
	aux_file.WriteString(string(buf_post)) // Armazena no novo arquivo

	Atualiza_indices("aux.txt", ind_arq)
	os.Rename("aux.txt", dados_arq)

	fmt.Println("Contato alterado.")

}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Vai atulizar o arquivo indice de acordo com arq de contatos
func Atualiza_indices(dados_arq, data_ind_nome string) {
	file, _ := os.Open(dados_arq)
	reader := bufio.NewReader(file)

	new_file_data, _ := os.OpenFile("aux_data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	new_file_ind, _ := os.OpenFile("aux_ind.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	defer new_file_data.Close()
	defer new_file_ind.Close()
	var x DataType

	for {
		nome, _, r_error := reader.ReadLine()
		tel, _, _ := reader.ReadLine()
		end, _, _ := reader.ReadLine()
		marca, _, _ := reader.ReadLine()
		reader.ReadLine()

		if r_error == io.EOF {
			break
		}

		end_of_file := FindIndEOF(new_file_data)
		x.nome = string(nome)
		x.indice = end_of_file
		SaveIndArq(x, new_file_ind)

		InsertWithtelAd(x, string(tel), string(end), string(marca), new_file_data)

	}
	os.Rename("aux_data.txt", dados_arq)
	os.Rename("aux_ind.txt", data_ind_nome)
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função responsavel por printat cada dado
func (contact DataType) Printar_Contato(file_data *os.File) {

	file_data.Seek(contact.indice, 0)

	reader := bufio.NewReader(file_data)
	nome, _, _ := reader.ReadLine()
	tel, _, _ := reader.ReadLine()
	end, _, _ := reader.ReadLine()
	marca, _, _ := reader.ReadLine()

	str_marca := string(marca)
	removed_marca, _ := strconv.Atoi(str_marca)

	if removed_marca == 1 {
		fmt.Printf("Nome: %s\n", nome)
		fmt.Printf("Endereco: %s\n", end)
		fmt.Printf("Telefone: %s\n\n", tel)
	}
	marca, _, _ = reader.ReadLine()
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que  carrega a arvore para memoria
func (tree *BTree) CarregaArv(ind_file, data_file string) {
	Ifile, err := os.Open(ind_file)
	reader := bufio.NewReader(Ifile)

	if err != nil {
		Ifile, err = os.OpenFile(ind_file, os.O_CREATE|os.O_WRONLY, 0644)
		Dfile, _ := os.OpenFile(data_file, os.O_CREATE|os.O_WRONLY, 0644)
		Ifile.Close()
		Dfile.Close()
		return
	}

	for {
		nome, _, r_error := reader.ReadLine()
		ind, _, _ := reader.ReadLine()
		ind_str := string(ind)

		if r_error == io.EOF {
			break
		}

		var x DataType

		x.nome = string(nome)
		ind_convert, _ := strconv.Atoi(ind_str)

		x.indice = int64(ind_convert)

		reader.ReadLine()

		tree.Insert(x)
	}
	Ifile.Close()
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que Esvazia a lixeira
func Esvazia_lixeira(dados_arq, ind_arq string) {
	file_data, _ := os.Open(dados_arq)
	reader := bufio.NewReader(file_data)

	new_file_data, _ := os.OpenFile("aux_data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	new_file_ind, _ := os.OpenFile("aux_ind.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file_data.Close()
	defer new_file_data.Close()
	defer new_file_ind.Close()
	var x DataType

	for {
		nome, _, r_error := reader.ReadLine()
		tel, _, _ := reader.ReadLine()
		end, _, _ := reader.ReadLine()
		marca, _, _ := reader.ReadLine()
		reader.ReadLine()

		if r_error == io.EOF {
			break
		}
		testing_marca, _ := strconv.Atoi(string(marca))

		if testing_marca == 1 { 
			end_of_file := FindIndEOF(new_file_data)
			x.nome = string(nome)
			x.indice = end_of_file
			SaveIndArq(x, new_file_ind)

			InsertWithtelAd(x, string(tel), string(end), string(marca), new_file_data)
		}
	}
	os.Rename("aux_data.txt", dados_arq) 
	os.Rename("aux_ind.txt", ind_arq)
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que faz a Impressão da árvore B em forma de árvore de diretório
func (node *BTreeNode) Print(indent string, last bool) {
	fmt.Print(indent)
	if last {
		fmt.Print("└─ ")
		indent += "    "
	} else {
		fmt.Print("├─ ")
		indent += "|   "
	}
	keys := make([]string, len(node.keys))
	fmt.Print("[")
	for i, key := range node.keys {
		keys[i] = fmt.Sprintf("%v", key)
	}
	fmt.Println(strings.Join(keys, "|"), "]")

	childCount := len(node.children)
	for i, child := range node.children {
		child.Print(indent, i == childCount-1)
	}
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que implementa a Divisão de um filho cheio
func (node *BTreeNode) splitChild(i int16) {
	child := node.children[i]
	newChild := InitNode(child.leaf)

	// Move as chaves e os filhos para o novo filho
	newChild.keys = append(newChild.keys, child.keys[t:]...)
	child.keys = child.keys[:t]
	if !child.leaf { // divide o nó em dois
		newChild.children = append(newChild.children, child.children[t:]...)
		child.children = child.children[:t]
	}

	// Insere o novo filho no nó
	node.children = append(node.children, nil)
	copy(node.children[i+2:], node.children[i+1:])
	node.children[i+1] = newChild

	// Move a chave correspondente para cima
	var aux DataType
	aux.nome = "0"
	aux.indice = 0
	node.keys = append(node.keys, aux)
	copy(node.keys[i+1:], node.keys[i:])
	node.keys[i] = child.keys[t-1]
	child.keys = child.keys[:t-1]
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que faz a Inserção de uma chave em um nó da árvore B
// Insert(key), key é a chave que será inserida
func (node *BTreeNode) Insert(key DataType) {
	if !node.leaf {
		// Encontra o filho apropriado para inserir a chave
		i := len(node.keys) - 1
		for i >= 0 && key.nome < node.keys[i].nome {
			i--
		}

		// Insere a chave no filho apropriado
		if len(node.children[i+1].keys) == 2*t-1 {
			node.splitChild(int16(i) + 1)
			if key.nome > node.keys[i+1].nome {
				i++
			}
		}
		node.children[i+1].Insert(key)
	} else {
		// Insere a chave no nó folha
		i := len(node.keys) - 1
		var aux DataType
		aux.nome = "0"
		aux.indice = 0
		node.keys = append(node.keys, aux)
		for i >= 0 && key.nome < node.keys[i].nome {
			node.keys[i+1] = node.keys[i]
			i--
		}
		node.keys[i+1] = key
	}
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que Dado o contato adiciona no arq de indices
func SaveIndArq(key DataType, file *os.File) {
	file.WriteString(string(key.nome))
	file.WriteString("\n")
	file.WriteString(string(strconv.Itoa(int(key.indice))))
	file.WriteString("\n\n")
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que Dado o contato adiciona  no arq de dados
func InsertWithtelAd(key DataType, tel, end, marca string, file *os.File) {
	file.WriteString(string(key.nome))
	file.WriteString("\n")
	file.WriteString(tel)
	file.WriteString("\n")
	file.WriteString(end)
	file.WriteString("\n")
	file.WriteString(marca)
	file.WriteString("\n\n")
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que acha o final do arquivo , em um arquivo já criado
func FindIndEOF(file *os.File) int64 { 
	end_of_file, _ := file.Seek(0, 2)

	return end_of_file
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que faz a Busca de chave
func (node *BTreeNode) ProcuraNome(key string) *DataType {
	i := 0
	for i < len(node.keys) && key > node.keys[i].nome {
		i++
	}

	if i < len(node.keys) && key == node.keys[i].nome {
		return &node.keys[i]
	} else if node.leaf {
		return nil
	} else {
		return node.children[i].ProcuraNome(key)
	}
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função de Inserir Chave na arvore
func (tree *BTree) Insert(key DataType) {
	root := tree.root
	if len(root.keys) == 2*t-1 {
		newRoot := InitNode(false)
		newRoot.children = append(newRoot.children, root)
		newRoot.splitChild(0)
		tree.root = newRoot
	}
	tree.root.Insert(key)
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// Função que faz a Busca de uma chave na árvore B
func (tree *BTree) ProcuraNome(key string) *DataType {
	return tree.root.ProcuraNome(key)
}

// ///////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////
// ////////////////////////// FUNÇÃO PRINCIPAL ////////////////////////////////
func main() {
	// Declaração de variáveis
	var tree *BTree
	var name DataType

	tree = Init()
	tree.CarregaArv(ArqInd, ArqDados)
	fmt.Println("======= Minha Agenda ===========")
	for {

		fmt.Println("========== Menu ================")
		fmt.Println()
		fmt.Println("1.  Incluir Contato             |")
		fmt.Println("2.  Listar Contatos             |")
		fmt.Println("3.  Excluir Contato             |")
		fmt.Println("4.  Restaura Contato da lixeira |")
		fmt.Println("5.  Esvaziar lixeira            |")
		fmt.Println("6.  Alterar Nome                |")
		fmt.Println("7.  Alterar Endereco            |")
		fmt.Println("8.  Alterar Telefone            |")
		fmt.Println("9.  Limpar Tela                 |")
		fmt.Println("10. Sair                        |")
		fmt.Println()
		fmt.Println("===============================")
		fmt.Println()
		fmt.Print("Escolha uma opção: ")

		var opcao int //leitura da opcao escolhida
		fmt.Scanln(&opcao)
		switch opcao {
		case 1:
			{

				var tel, end string

				fmt.Println()
				fmt.Println("Entre com os Dados")
				fmt.Println(" ")
				fmt.Print("Nome : ")
				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.TrimRight(input, "\n")
				name.nome = input
				
				// Verifica se está no tamanho permetido
				if len(name.nome) > MaxNome {
					fmt.Println(" ")
					fmt.Printf("Erro: O nome excede o tamanho máximo permitido (%d Caracteres)\n", MaxNome)
					for {
						fmt.Println(" ")
						fmt.Printf("Insira um nome valido !!! \n")
						fmt.Print("Nome : ")
						input, _ := reader.ReadString('\n')
						input = strings.TrimRight(input, "\n")
						name.nome = input
						//Obriga o usuario entrar com um valor valido
						if len(name.nome) > MaxNome {
							fmt.Printf("Erro: O nome excede o tamanho máximo permitido (%d)\n", MaxNome)
						} else {
							fmt.Println(" ")
							break
						}
					}
				}

				// Verifica se está no tamanho permetido
				fmt.Print("Endereco : ")
				input, _ = reader.ReadString('\n')
				input = strings.TrimRight(input, "\n")
				end = input

				if len(end) > MaxEndereco {
					fmt.Printf("Erro: O endereço excede o tamanho máximo permitido (%d Caracteres)\n", MaxEndereco)
					for {
						fmt.Println(" ")
						fmt.Printf("Insira um endereço valido !!! \n")
						fmt.Print("Endereço : ")
						reader_rem := bufio.NewReader(os.Stdin)
						remov, _ := reader_rem.ReadString('\n')
						remov = strings.TrimRight(remov, "\n")
						input, _ := reader.ReadString('\n')
						input = strings.TrimRight(input, "\n")
						end = input
						//Obriga o usuario entrar com um valor valido
						if len(end) > MaxEndereco {
							fmt.Printf("Erro: O endereço excede o tamanho máximo permitido (%d)\n", MaxEndereco)
						} else {
							fmt.Println(" ")
							break
						}
					}
				}

				fmt.Print("Numero Telefonico : ")
				input, _ = reader.ReadString('\n')
				input = strings.TrimRight(input, "\n")
				tel = input
				fmt.Println(" ")

				// Verifica se está no tamanho permetido
				if len(tel) > MaxTelefone {
					fmt.Printf("Erro: O telefone excede o tamanho máximo permitido (%d Caracteres)\n", MaxTelefone)
					for {
						fmt.Println(" ")
						fmt.Printf("Insira um Telefone valido !!! \n")
						fmt.Print("Telefone : ")
						input, _ := reader.ReadString('\n')
						input = strings.TrimRight(input, "\n")
						tel = input
						//Obriga o usuario entrar com um valor valido
						if len(tel) > MaxTelefone {
							fmt.Printf("Erro: O telefone excede o tamanho máximo permitido (%d)\n", MaxTelefone)
						} else {
							fmt.Println(" ")
							break
						}
					}
				}

				data_file, errD := os.OpenFile(ArqDados, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				ind_file, errI := os.OpenFile(ArqInd, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

				if errD != nil || errI != nil {
					fmt.Println("Arquivo inexistente ")
				}

				ind := FindIndEOF(data_file)
				// Insererindo elemento no arquivo
				InsertWithtelAd(name, tel, end, "1", data_file)
				name.indice = ind
				SaveIndArq(name, ind_file)

				tree.Insert(name)
				fmt.Println("Contatos salvos ;-)")
				//fmt.Scanln()
				ind_file.Close()
				data_file.Close()
				fmt.Println(" ")
			}

		case 2:
			{
				fmt.Println(" ")
				fmt.Println("Todos Contatos Armazenados")
				fmt.Println(" ")
				tree.root.Listar_Ord(ArqDados)

			}

		case 3:
			{
				fmt.Println(" ")
				fmt.Println("Entre com o contato a ser removido")
				reader_rem := bufio.NewReader(os.Stdin)
				remov, _ := reader_rem.ReadString('\n')
				remov = strings.TrimRight(remov, "\n")
				tree.Remove_Contato(remov, ArqDados)

			}
		case 4:
			{
				fmt.Println(" ")
				fmt.Println("Recuperar Lixeira")
				RestauraLixeira(ArqDados)
				fmt.Println("Elementos Recuperados.")
				fmt.Println(" ")
			}
		case 5:
			{
				fmt.Println(" ")
				fmt.Println("Esvaziando lixeira ...")
				Esvazia_lixeira(ArqDados, ArqInd)
				tree = Init()
				tree.CarregaArv(ArqInd, ArqDados)
				fmt.Println("Lixeira vazia !")
				fmt.Println(" ")
			}

		case 6:
			{
				fmt.Println(" ")
				fmt.Println("Entre com o Nome Contato a ser alterado : ")
				reader_mod := bufio.NewReader(os.Stdin)
				nome, _ := reader_mod.ReadString('\n')
				nome = strings.TrimRight(nome, "\n")

				fmt.Println("Novo Nome: ")
				novo_nome, _ := reader_mod.ReadString('\n')
				novo_nome = strings.TrimRight(novo_nome, "\n")

				tree.Altera_nome(nome, novo_nome, ArqDados, ArqInd)
				tree = Init()
				tree.CarregaArv(ArqInd, ArqDados)
				fmt.Println(" ")
			}

		case 7:
			{
				fmt.Println(" ")
				fmt.Println("Entre com o Nome Contato a ser alterado : ")

				reader_mod := bufio.NewReader(os.Stdin)
				nome, _ := reader_mod.ReadString('\n')
				nome = strings.TrimRight(nome, "\n")

				fmt.Println("Novo endereco : ")
				novo_end, _ := reader_mod.ReadString('\n')
				novo_end = strings.TrimRight(novo_end, "\n")

				tree.Altera_end(nome, novo_end, ArqDados, ArqInd)
				tree = Init()
				tree.CarregaArv(ArqInd, ArqDados)
				fmt.Println(" ")
			}
		case 8:
			{
				fmt.Println(" ")
				fmt.Println("Entre com o Nome Contato a ser alterado : ")

				reader_mod := bufio.NewReader(os.Stdin)
				nome, _ := reader_mod.ReadString('\n')
				nome = strings.TrimRight(nome, "\n")

				fmt.Println("Novo Telefone : ")

				novo_tel, _ := reader_mod.ReadString('\n')
				novo_tel = strings.TrimRight(novo_tel, "\n")

				tree.Alterar_tel(nome, novo_tel, ArqDados, ArqInd)
				tree = Init()
				tree.CarregaArv(ArqInd, ArqDados)
				fmt.Println(" ")
			}
		case 9:
			{
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				cmd.Run()
			}

		case 10:
			{
				fmt.Println(" ")
				fmt.Println("Saindo...") //fecha o programa
				return
			}
		default:
			{
				fmt.Println(" ")
				fmt.Println("Operação Inválida !")
				fmt.Println("Operações Validas são do 1 ao 10 !!!")
				fmt.Println(" ")
			}
		}
	}
}
