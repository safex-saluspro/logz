# **Contribuindo para Logz**

Obrigado por se interessar em contribuir para o **Logz**! Estamos animados por ter você como parte da nossa comunidade. Este guia ajudará você a começar e contribuir de forma eficiente com o projeto.

---

## **Como Contribuir**

Existem várias maneiras de contribuir com o Logz:

1. **Relatar Problemas (Issues)**
    - Encontre bugs ou problemas no código? Abra uma _issue_ detalhando o problema.
    - Inclua o máximo de informações possíveis: passos para reproduzir o problema, logs, versão do Go utilizada, etc.

2. **Sugerir Melhorias**
    - Tem uma ideia para melhorar o projeto? Compartilhe sua sugestão abrindo uma _issue_ com a tag `enhancement`.

3. **Enviar Pull Requests**
    - Quer corrigir um bug ou implementar algo novo? Envie um _pull request_ com suas alterações.

4. **Testar e Revisar Código**
    - Ajude a revisar _pull requests_ de outros contribuidores.
    - Execute os testes existentes e valide se as alterações propostas mantêm o sistema funcional.

---

## **Primeiros Passos**

### 1. **Clone o Repositório**
```bash
git clone https://github.com/<seu-usuario>/logz.git
cd logz
```

### 2. **Configure o Ambiente**
Certifique-se de ter o Go instalado:
- [Download Go](https://go.dev/dl/)

### 3. **Instale Dependências**
```bash
# Baixe os pacotes necessários
go mod download
```

### 4. **Execute os Testes**
Antes de fazer mudanças, execute os testes existentes:
```bash
go test ./...
```

---

## **Criando um Pull Request**

### **1. Faça um Fork do Repositório**
Crie um fork do projeto para seu próprio GitHub.

### **2. Crie uma Branch Nova**
```bash
git checkout -b sua-feature
```

### **3. Faça as Alterações**
Certifique-se de seguir as convenções de código e boas práticas do projeto.

### **4. Adicione Testes (se aplicável)**
Inclua casos de teste para validar a funcionalidade adicionada.

### **5. Execute os Testes**
Certifique-se de que todas as alterações e testes estão funcionando:
```bash
go test ./...
```

### **6. Commit e Push**
```bash
git add .
git commit -m "Descrição breve da alteração"
git push origin sua-feature
```

### **7. Abra o Pull Request**
Vá até o repositório original no GitHub e abra um _pull request_ explicando suas mudanças.

---

## **Padrões de Código**

### **Estilo de Código**
Este projeto segue as convenções de código do Go. Algumas recomendações:
- Use `gofmt` para formatar o código:
```bash
gofmt -w .
```

- Nomeie variáveis e funções de forma clara e descritiva.
- Divida funções longas em partes menores sempre que possível.

### **Commits**
Os commits devem ser claros e descritivos. Exemplos:
- `fix: corrigir bug na lógica de notificação`
- `feat: adicionar suporte ao notifier via Slack`

---

## **Boas Práticas**

1. **Seja Respeitoso e Acolhedor**  
   Este é um projeto de código aberto para todos. Respeite outros contribuidores e colabore de forma construtiva.

2. **Documente Suas Alterações**  
   Atualize o `README.md` ou a documentação, se necessário, para incluir suas mudanças.

3. **Adicione Testes Quando Possível**  
   Assegure-se de que qualquer nova funcionalidade seja acompanhada por testes.

4. **Seja Claro nos Relatórios de Problemas**  
   Ao abrir uma _issue_, seja detalhado e forneça o máximo de contexto.

---

## **Onde Obter Ajuda**

Se você precisar de assistência, sinta-se à vontade para:
- Abrir uma _issue_ com a tag `question`.
- Entrar em contato comigo através do e-mail ou LinkedIn listados no `README.md`.

---

## **Nosso Compromisso**

Nos comprometemos a revisar pull requests e issues o mais rápido possível. Valoramos sua contribuição e agradecemos pelo tempo dedicado ao projeto!
