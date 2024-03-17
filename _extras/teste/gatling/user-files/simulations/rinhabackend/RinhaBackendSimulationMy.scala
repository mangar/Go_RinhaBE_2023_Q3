import java.time.LocalDate

import scala.concurrent.duration._
import scala.util.Random

import io.gatling.core.Predef._
import io.gatling.http.Predef._


class RinhaBackendSimulationMy extends Simulation {

  def randomNome(): String = {
    val nomes = Array("Lucas", "Maria", "João", "Ana", "Pedro", "Julia", "Marcos", "Luana", "Marcio", "Alexandre", "Rafael", "Patricia", "Thomaz")
    val sobrenomes = Array("Silva", "Santos", "Oliveira", "Souza", "Rodrigues", "Fernandes", "Pereira", "Lima")

    val random = new Random()
    val nomeAleatorio = nomes(random.nextInt(nomes.length))
    val sobrenomeAleatorio = sobrenomes(random.nextInt(sobrenomes.length))
    s"\"${nomeAleatorio} ${sobrenomeAleatorio}\""
  }
  def randomApelido() = s"\"${Random.alphanumeric.take(10).mkString}\""
  def randomDataNascimento(): String = {    
    val minAge = 18; val maxAge = 80
    val random = new Random()
    val anoAtual = LocalDate.now().getYear
    val anoAleatorio = anoAtual - minAge - random.nextInt(maxAge - minAge + 1)
    val mesAleatorio = random.nextInt(12) + 1
    val diaAleatorio = mesAleatorio match {
      case 2 if anoAleatorio % 4 == 0 => random.nextInt(29) + 1
      case 2 => random.nextInt(28) + 1
      case 4 | 6 | 9 | 11 => random.nextInt(30) + 1
      case _ => random.nextInt(31) + 1
    }

    f"\"$anoAleatorio%04d-$mesAleatorio%02d-$diaAleatorio%02d\""
  }  
  def stacks(): Array[String] = Array("Scala", "Java", "Kotlin", "Python", "Go", "Rust", "JavaScript", "TypeScript", "C#", 
      "Oracle", "SQLServer", "MySQL", "MongoDB", "CouchDB", "Redis",  
      "AWS", "GCP", "Azure")
  val feederStack = stacks().map(stack => Map("stack" -> stack)).circular  
  def randomStacks(): String = {
    val arrayOriginal = stacks()
    val random = new Random()
    val novoTamanho = random.nextInt(10) + 1  // Gerar tamanho aleatório entre 1 e 10
    val novo = Array.fill(novoTamanho)(arrayOriginal(random.nextInt(arrayOriginal.length)))
    novo.mkString("[\"", "\", \"", "\"]")
  }


  def randomApelidoInvalidoNull() = Seq(null, "")(Random.between(0, 1 + 1))
  def randomNomeInvalidoNull() = Seq(null, "")(Random.between(0, 1 + 1))


  val httpProtocol = http
    .baseUrl("http://localhost:3000").userAgentHeader("Agente do Caos - 2023")

  val criacaoEConsultaPessoas_Validas = scenario("Criação E Talvez Consulta de Pessoas - Validas")
    .exec{ s => 
      val apelido = randomApelido()
      val nome = randomNome()
      val dataNascimento = randomDataNascimento()
      val stack = randomStacks()
      val payload = s"""{"apelido" : ${apelido}, "nome" : ${nome}, "nascimento" : ${dataNascimento}, "stack" : ${stack}}"""
      val session = s.setAll(Map("apelido" -> apelido, "payload" -> payload))      
      session
    }
    .exec(
      http("Criação Valida")
      .post(s => "/pessoas")
        .header("content-type", "application/json")
        .body(StringBody(s => s("payload").as[String]))
      
        // 201 pros casos de sucesso :)
        // 422 pra requests inválidos :|
        // 400 pra requests bosta tipo data errada, tipos errados, etc. :(
        .check(
          status.in(201, 422, 400),
          status.saveAs("httpStatus"))

        .checkIf(s => s("httpStatus").as[String] == "201") {
          header("Location").saveAs("location")
        }
    )
    .pause(1.milliseconds, 30.milliseconds)
    .doIf(s => s.contains("location")) {
      exec(
        http("Criação Valida > Consulta")
        .get("${location}")
          .header("content-type", "application/json")
          .check(
            status.in(200)
          )
      )
    }



  val criacaoEConsultaPessoas_NOTValidas = scenario("Criação E Talvez Consulta de Pessoas - NOT Validas")
    .exec{ s => 
      val apelido = randomApelidoInvalidoNull()
      val nome = randomNomeInvalidoNull()
      val dataNascimento = randomDataNascimento()
      val stack = randomStacks()
      val payload = s"""{"apelido" : ${apelido}, "nome" : ${nome}, "nascimento" : ${dataNascimento}, "stack" : ${stack}}"""
      val session = s.setAll(Map("apelido" -> apelido, "payload" -> payload))      
      session
    }
    .exec(
      http("Criação IN Valida")
      .post(s => "/pessoas")
        .header("content-type", "application/json")
        .body(StringBody(s => s("payload").as[String]))

        // 422 pra requests inválidos :|
        .check(status.in(422))
    )
    








  val buscaPessoas = scenario("Busca Válida de Pessoas")
    // .exec{ s => 
    //   val feederStack = stacks().map(stack => Map("stack" -> stack)).circular  
    //   val session = s.setAll(Map("feederStack" -> feederStack))
    //   session    
    // }    
    .feed(feederStack)
      .exec(
        http("Busca por Stack")
          .get("/pessoa?t=#{stack}")  // Usando o dado injetado na URL
          .header("content-type", "application/json")
          .check(status.in(200, 404))
    )

    // .feed(tsv("termos-busca.tsv").circular())
    // .exec(
    //   http("busca válida")
    //   .get("/pessoas?t=#{t}")
    //   // qq resposta na faixa 2XX tá safe
    // )

  // val buscaInvalidaPessoas = scenario("Busca Inválida de Pessoas")
  //   .exec(
  //     http("busca inválida")
  //     .get("/pessoas")
  //     // 400 - bad request se não passar 't' como query string
  //     .check(status.is(400))
  //   )

  setUp(
    criacaoEConsultaPessoas_Validas.inject(
      // constantUsersPerSec(2).during(10.seconds), // warm up
      // constantUsersPerSec(5).during(15.seconds).randomized, // are you ready?
      rampUsersPerSec(1).to(220).during(2.minutes),
      constantUsersPerSec(220).during(2.minutes)
    )
    // ,criacaoEConsultaPessoas_NOTValidas.inject(
    //   constantUsersPerSec(2).during(10.seconds), // warm up
    //   constantUsersPerSec(5).during(15.seconds).randomized, // are you ready?
    //   rampUsersPerSec(6).to(300).during(1.minute) // lezzz go!!!
    // )    
    // ,buscaPessoas.inject(
    //   constantUsersPerSec(2).during(10.seconds), // warm up

    //   // rampUsersPerSec(6).to(50).during(3.minutes) // lezzz go!!!
    // )


    // buscaInvalidaPessoas.inject(
    //   constantUsersPerSec(2).during(25.seconds), // warm up

    //   rampUsersPerSec(6).to(20).during(3.minutes) // lezzz go!!!
    // )
  ).protocols(httpProtocol)
}
