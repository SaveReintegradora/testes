Feature: Livros
  Para garantir que a API de livros funciona corretamente
  Como um usuário da API
  Eu quero criar e listar livros

  Scenario: Criar um novo livro
    Given que eu defino a API key para "minha-chave-secreta"
    And que eu crio um livro com título "Livro Teste" e autor "Autor"
    When eu envio uma requisição POST para o endpoint "/books"
    Then o status da resposta deve ser 201
    And o corpo da resposta deve conter "id"
    And o corpo da resposta deve conter "title"
    And o corpo da resposta deve conter "author"

  Scenario: Listar livros
    Given que eu defino a API key para "minha-chave-secreta"
    When eu envio uma requisição GET para o endpoint "/books"
    Then o status da resposta deve ser 200
    And o corpo da resposta deve conter "title"
    And o corpo da resposta deve conter "author"
