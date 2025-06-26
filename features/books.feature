# language: pt
Funcionalidade: Gerenciamento de Livros

  Contexto:
    Dado que estou autenticado com a API Key "minha-chave-secreta"

  Cenário: Listar todos os livros
    Quando faço uma requisição GET para "/books"
    Então a resposta deve ter status 200
    E a resposta deve conter uma lista de livros

  Cenário: Criar um novo livro
    Quando faço uma requisição POST para "/books" com o corpo:
      """
      {"title": "Livro Teste", "author": "Autor"}
      """
    Então a resposta deve ter status 201
    E a resposta deve conter o livro criado com título "Livro Teste"

  Cenário: Buscar livro por ID existente
    Dado que existe um livro com ID "algum"
    Quando faço uma requisição GET para "/books/ultimo"
    Então a resposta deve ter status 200
    E a resposta deve conter o livro com ID "${ultimo_id}"

  Cenário: Atualizar um livro existente
    Dado que existe um livro com ID "algum"
    Quando faço uma requisição PUT para "/books/ultimo" com o corpo:
      """
      {"title": "Novo Título"}
      """
    Então a resposta deve ter status 200
    E a resposta deve conter o livro com título "Novo Título"

  Cenário: Remover um livro existente
    Dado que existe um livro com ID "algum"
    Quando faço uma requisição DELETE para "/books/ultimo"
    Então a resposta deve ter status 204

  Cenário: Buscar livro por ID inválido
    Quando faço uma requisição GET para "/books/nao-uuid"
    Então a resposta deve ter status 400
    E a resposta deve conter o campo "error" com valor "ID inválido"
