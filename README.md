
# Ze Code Challenge API

## Visão geral
Este projeto é baseado no desafio de backend do Zé Delivery.

É uma API Rest que dado um parceiro e a area de cobertura de sua entrega implemente as seguintes funcionalidades e requisitos técnicos:

### 1. Criar um parceiro:
Salvar no banco de dados todas as seguintes informações representadas por este JSON junto com as regras subsequentes:

### 2. Carregar parceiro pelo `id`:
Retornar um parceiro específico baseado no seu campo `id` com todos os campos apresentados acima.

### 3. Buscar parceiro:
Dada uma localização pelo usuário da API (coordenadas `long` e `lat`), procure o parceiro que esteja **mais próximo** e **que cuja área de cobertura inclua** a localização.

```JSON
  [
    {
      "id":1,
      "tradingName":"Adega da Cerveja - Pinheiros",
      "ownerName":"Ze da Silva",
      "document":"1432132123891/0001",
      "coverageArea":{
        "type":"MultiPolygon",
        "coordinates":[
          [[[30,20],[45,40],[10,40],[30,20]]],
          [[[15,5],[40,10],[10,20],[5,10],[15,5]]]]
          },
      "address":{
        "type":"Point",
        "coordinates":[-46.57421,-21.785741]
        }
    },
  ]
```
1. O campo `address` (endereço em inglês) segue o formato `GeoJSON Point` (https://en.wikipedia.org/wiki/GeoJSON);
2. o campo `coverageArea` (área de cobertura em inglês) segue o formato `GeoJSON MultiPolygon` (https://en.wikipedia.org/wiki/GeoJSON);
3. O campo `document` deve ser único entre os parceiros;
4. O campo `id` deve ser único entre os parceiros, mas não necessariamente um número inteiro;
## Documentação da API

###  Parceiro
| Campo        | Descrição                              | Tipo    |
|--------------|----------------------------------------|---------|
| id           | Identificador único do parceiro        | int     |
| tradingName  | Nome comercial do parceiro             | string  |
| document     | Documento de identificação do parceiro | string  |
| coverageArea | Area de cobertura de entregas          | GeoJSON MultiPolygon|
| address      | Endereço do parceiro                   | GeoJSON Point|


#### Retorna todos os parceiros

```http
  GET /partners
```

```JSON
  [
    {"id":1,
    "tradingName":"Adega da Cerveja - Pinheiros",
    "ownerName":"Ze da Silva",
    "document":"1432132123891/0001",
    "coverageArea":{
      "type":"MultiPolygon",
      "coordinates":[
        [[[30,20],[45,40],[10,40],[30,20]]],
        [[[15,5],[40,10],[10,20],[5,10],[15,5]]]]
        },
    "address":{
      "type":"Point",
      "coordinates":[-46.57421,-21.785741]
      }
    },
  ]
```

#### Retorna o parceiro pelo ID

```http
  GET /partners/{$id}
```

```JSON
{
  "id":1,
  "tradingName":"Adega da Cerveja - Pinheiros",
  "ownerName":"Ze da Silva",
  "document":"1432132123891/0001",
  "coverageArea":{
    "type":"MultiPolygon",
    "coordinates":[
      [[[30,20],[45,40],[10,40],[30,20]]],
      [[[15,5],[40,10],[10,20],[5,10],[15,5]]]]
      },
  "address":{
    "type":"Point",
    "coordinates":[-46.57421,-21.785741]
  }
},
```

#### Registra um novo parceiro


```http
  POST /partners
```

```JSON
{
  "tradingName": "Adega do Opalão - Mooca",
  "ownerName": "Raul Opala",
  "document": "2032132123891/2001",
  "coverageArea": { 
    "type": "MultiPolygon", 
    "coordinates": [
      [[[12, 10], [16, 25], [37, 25], [46, 19]]] 
    ]
  },
  "address": { 
    "type": "Point",
    "coordinates": [-10, -6]
  }
}
```


#### Busca o parceiro mais proximo cuja a área de cobertura inclui a localização

| Campo       | Descrição                   | Tipo    |
|-------------|-----------------------------|---------|
| x           | Latitude                    | float64 |
| y           | Longitude                   | float64 |
| maxDistance | Distancia máxima do usuário | float64 |


```http
  POST /partners/near
```

```JSON
{
    "X": 31,
    "Y": 12,
    "maxDistance": 90
}
```


## Implementação

O projeto foi desenvolvido com [Go](https://go.dev) e o pacote [Gin](https://github.com/gin-gonic/gin) uma framework em Go 40x mais rápida.

```bash
  go run .
```

