package models

import (
	"fmt"
	"strconv"

	"github.com/astaxie/beego/orm"
)

type Paragrafo struct {
	Id     int
	Numero int
	Texto  string
}

type Articulo struct {
	Id         int
	Numero     int
	Texto      string
	Paragrafos []Paragrafo
}

type ResolucionCompleta struct {
	Vinculacion   ResolucionVinculacionDocente
	Consideracion string
	Preambulo     string
	Vigencia      int
	Numero        string
	Id            int
	Articulos     []Articulo
	Titulo        string
}

func GetOneResolucionCompleta(idResolucion string) (resolucion ResolucionCompleta) {
	o := orm.NewOrm()
	var temp []Resolucion
	idRes, _ := strconv.Atoi(idResolucion)

	_, err := o.QueryTable("resolucion").Filter("id_resolucion", idRes).All(&temp)
	if err == nil {
		fmt.Println("Consulta exitosa")
	}

	resolucionCompleta := ResolucionCompleta{Id: temp[0].Id, Consideracion: temp[0].ConsideracionResolucion, Preambulo: temp[0].PreambuloResolucion, Vigencia: temp[0].Vigencia, Numero: temp[0].NumeroResolucion, Titulo: temp[0].Titulo}

	var arts []ComponenteResolucion
	_, err2 := o.QueryTable("componente_resolucion").Filter("resolucion_id", idRes).Filter("tipo_componente", "Articulo").OrderBy("numero").All(&arts)
	if err2 == nil {
		fmt.Println("Consulta exitosa")
	}

	var articulos []Articulo

	for _, art := range arts {
		articulo := Articulo{Id: art.Id, Numero: art.Numero, Texto: art.Texto}

		var pars []ComponenteResolucion
		_, err3 := o.QueryTable("componente_resolucion").Filter("resolucion_id", idRes).Filter("tipo_componente", "Paragrafo").Filter("componente_padre", articulo.Id).OrderBy("numero").All(&pars)
		if err3 == nil {
			fmt.Println("Consulta exitosa")
		}

		var paragrafos []Paragrafo

		for _, par := range pars {
			paragrafo := Paragrafo{Id: par.Id, Numero: par.Numero, Texto: par.Texto}
			paragrafos = append(paragrafos, paragrafo)
		}

		articulo.Paragrafos = paragrafos

		articulos = append(articulos, articulo)
	}
	resolucionCompleta.Articulos = articulos
	return resolucionCompleta
}

func UpdateResolucionCompletaById(m *ResolucionCompleta) (err error) {
	o := orm.NewOrm()
	v := Resolucion{Id: m.Id}
	if err = o.Read(&v); err == nil {
		v.NumeroResolucion = m.Numero
		v.Titulo = m.Titulo
		_, err = o.Update(&v)
	} else {
		return
	}
	idResolucionStr := strconv.Itoa(m.Id)
	r := m.Vinculacion
	fmt.Println(r.Id)
	a := ResolucionVinculacionDocente{Id: r.Id}
	if err = o.Read(&a); err == nil {
		_, err = o.Update(&r)
	} else {
		return
	}
	if err = o.Read(&v); err == nil {
		v.ConsideracionResolucion = m.Consideracion
		v.PreambuloResolucion = m.Preambulo
		v.NumeroResolucion = m.Numero
		fmt.Println(v)
		if err := UpdateResolucionById(&v); err != nil {
		}

		resolucionCompleta := GetOneResolucionCompleta(idResolucionStr)

		for _, articulo := range resolucionCompleta.Articulos {
			if articulo.Paragrafos != nil {
				for _, paragrafo := range articulo.Paragrafos {
					if err := DeleteComponenteResolucion(paragrafo.Id); err != nil {
					}
				}
			}
			if err := DeleteComponenteResolucion(articulo.Id); err != nil {
			}
		}

		for indexArticulo, articulo := range m.Articulos {
			componenteArticulo := ComponenteResolucion{ResolucionId: &Resolucion{Id: m.Id}, Texto: articulo.Texto, Numero: indexArticulo + 1, TipoComponente: "Articulo"}
			if _, err := AddComponenteResolucion(&componenteArticulo); err == nil {
				if articulo.Paragrafos != nil {
					for indexParagrafo, paragrafo := range articulo.Paragrafos {
						componenteParagrafo := ComponenteResolucion{ResolucionId: &Resolucion{Id: m.Id}, Texto: paragrafo.Texto, Numero: indexParagrafo + 1, TipoComponente: "Paragrafo", ComponentePadre: &ComponenteResolucion{Id: componenteArticulo.Id}}
						if _, err := AddComponenteResolucion(&componenteParagrafo); err == nil {

						}
					}
				}
			}
		}
	}
	return
}

func GetTemplateResolucion(dedicacion, nivel, periodo, tipo string) (res ResolucionCompleta) {
	var resolucion ResolucionCompleta
	var articulos []Articulo
	var articulo Articulo
	var paragrafo Paragrafo
	//var vigencia, _, _ = time.Now().Date()
	//var accion string
	//var periodoStr string
	//var nombreDedicacion string
	/*
		switch periodo {
		case "1":
			periodoStr = "primer"
		case "2":
			periodoStr = "segundo"
		case "3":
			periodoStr = "tercer"
		}
	*/

	switch dedicacion {
	case "HCP":
		if nivel == "PREGRADO" {
			resolucion = ResolucionCompleta{Consideracion: "Que el artículo tercero del Decreto 1279 de junio 19 de 2002, mediante el cual se establece el régimen salarial y prestacional de los docentes de las universidades estatales, señala que “[l]os profesores ocasionales no son empleados públicos docentes de régimen especial ni pertenecen a la carrera profesoral y, por consiguiente, sus condiciones salariales y prestacionales no están regidas por el presente Decreto”, precisando que, “[n]o obstante, su vinculación se hace conforme a las reglas que define cada Universidad, con sujeción a lo dispuesto por la ley 30 de 1992 y demás disposiciones constitucionales y legales vigentes”.\n\nQue la Universidad Distrital Francisco José de Caldas puede vincular docentes de vinculación especial, en las modalidades de hora cátedra (HC), medio tiempo ocasional (MTO) y tiempo completo ocasional (TCO), Honorarios hora catedra, en virtud de lo consagrado en el artículo 13 del Acuerdo 011 de noviembre 15 de 2002 (Estatuto Docente), expedido por el Consejo Superior Universitario, a término fijo, por períodos académicos.\n\nQue mediante Resolución Nro. 001 del 15 de febrero de 2012, proferida por la Vicerrectoría Académica, se establece el procedimiento para la selección y vinculación a la Universidad, de docentes de vinculación especial.\n\nQue los servicios de los docentes de vinculación especial, deberán ser reconocidos en los términos del inciso segundo del artículo 74 de la Ley 30 de 1992, esto es, “mediante resolución”\n\nQue la Corte Constitucional en Sentencia C-006 de 1996, a través de la cual se declararon inexequibles algunos apartes del artículo 73 de la Ley 30 de 1992, aclaró, entre otras cosas, que, la vinculación de los docentes de hora cátedra estará mediada por su transitoriedad o temporalidad, de lo cual se deriva la proporcionalidad prestacional al tiempo laborado.\n\nQue conforme al concepto de fecha diciembre 18 de 2018, emitido por el Departamento Administrativo de la Función Pública, “en el caso de la Universidad Distrital Francisco José de Caldas, debe acogerse  a lo dispuesto en sus estatutos internos o en el acuerdo en el cual defina sus políticas y criterios generales para regular los procesos de selección, vinculación y contratación de los docentes ocasionales y catedráticos, en especial en lo relacionado con el reconocimiento y pago de prestaciones sociales y remuneración…”.\n\nQue conforme al parágrafo 1º del artículo quinto de la Resolución Nro. 001 de 2012, “Para efectos de pago y liquidación el mes comprenderá (4) semanas o (30) días laborales en Proyectos Académicos de Pregrado y en los espacios académicos de Posgrado se pagará de acuerdo al número de horas efectivamente dictadas en cada período”.\n\nQue el artículo 2 del Decreto 310 de 27 febrero 2020, emanado por el Presidente de la República de Colombia, establece, “ (…) a  partir del 1° de enero del 2020, fijase el valor del punto para empleados públicos docentes a quienes se les aplica el Decreto 1279 de 2002 y demás disposiciones que lo modifiquen o adicionen en catorce mil novecientos treinta y ocho pesos ($14.938) moneda corriente (…)”.\n\nQue, mediante Resolución 137 de abril 3 de 2020, emitida por Rectoría de la Universidad Distrital Francisco José de Caldas, se determinó “acoger el artículo 2 del Decreto 310 de 27 de febrero de 2020, en lo pertinente, única y expresamente al valor del punto para los docentes de Vinculación Especial hora cátedra, honorarios, Medio Tiempo Ocasional y Tiempo Completo Ocasional en pregrado, el valor del punto en CATORCE MIL NOVECIENTOS TREINTA Y OCHO PESOS ($14.938) MONEDA CORRIENTE”.\n\nQue en materia de Seguridad y Salud en el Trabajo SST para docentes ocasionales de la Universidad Distrital se deben informar las obligaciones específicas se encuentra contemplado en el Decreto 1072 de 2015 artículo 2.2.4.2.2.16. y la Resolución de Rectoría No. 624 de 2017, el docente ocasional, hora catedra y honorarios debe cumplir con las normas del sistema general de riesgos laborales.\n\nQue el Consejo Académico en uso de las funciones estatutarias establecidas en el Artículo 18, literales b.) y e.), del Estatuto General de la Universidad -Acuerdo 003 de abril 08 de 1997, expedido por el Consejo Superior Universitario-, mediante Resolución N°  053 de julio 23 del 2020 aprobó “el Calendario Académico especial para el segundo semestre del año 2020, para los programas académicos de pregrado bajo la modalidad virtual y/o alternancia, en el marco de las medidas de emergencia sanitaria y de aislamiento social decretado por el Gobierno Nacional, por cuenta de la pandemia Covid-19” \n\nQue la Resolución N°. 053 del 23 de julio del 2020, el Consejo Académico determinó en el artículo primero que las clases se inician el día 28 de septiembre de 2020 y se suspenden a partir del 20 de diciembre de 2020, cuando se entra en un periodo de receso hasta el 24 de enero de 2021. \n\nQue la Rectoría de la Universidad Distrital Francisco José de Caldas en uso de la autonomía universitaria y ante la problemática de salud pública, por medio de la Resolución 132 de marzo 19 de 2020, adoptó medidas transitorias y excepcionales de orden interno, de acuerdo con los lineamientos del orden Nacional y Distrital, para afrontar debidamente la situación de emergencia epidemiológica, evitando el contagio y garantizando la prestación del servicio a cargo de la Universidad.\n\nQue entre las medidas transitorias y excepcionales adoptadas en la Resolución de Rectoría 132 de marzo 19 de 2020, se optó por la “suspensión de las clases presenciales en todos los programas pregrado y posgrado de conformidad con el comunicado del 15 de marzo de 2020, a partir de las seis (6) de la mañana del 19 de marzo de 2020 y hasta nuevas instrucciones, según la actual coyuntura de salud pública, así como de emergencia económica y social”.\n\nQue con base en la Directiva Ministerial Nº 04 del 22 de marzo de 2020 del Ministerio de Educación Nacional, y con el propósito de dar continuidad durante el periodo de emergencia sanitaria a los programas académicos con registro calificado con que cuenta la Universidad Distrital Francisco José de Caldas en modalidad presencial, de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\nQue la misma Directiva Ministerial prevé que una vez finalizada la emergencia sanitaria, los desarrollos curriculares de los programas académicos con registro calificado en modalidad presencial deberán desarrollarse como se venía realizando habitualmente de acuerdo con las características propias de dicha modalidad.\n\nQue conforme lo dispuesto en el artículo 15 del Decreto Legislativo 491 de 2020, expedido por el Presidente de la República de Colombia, tiene por objeto que las autoridades cumplan con la finalidad de proteger y garantizar los derechos y libertades de las personas, la primacía de los intereses generales, la sujeción de las autoridades a la Constitución y demás preceptos del ordenamiento jurídico, el cumplimiento de los fines y principios estatales, el funcionamiento eficiente y democrático de la administración y la observancia de los deberes del Estado y de los particulares, durante el período de aislamiento preventivo obligatorio, las universidades públicas dispondrán de las medidas necesarias para que los docentes ocasionales o de hora cátedra cumplan sus funciones mediante la modalidad de \"trabajo en casa\", entre otras, haciendo uso de las tecnologías de la información y las comunicaciones.\n\nQue se hace necesario garantizar la vinculación de docentes hora cátedra que realicen las actividades lectivas señaladas en el calendario académico y que no son suplidas por los docentes de carrera de la Universidad Distrital.\n\nQue, para efectos presupuestales, el presente acto administrativo se ejecutará con cargo a los recursos de que trata los Certificados de Disponibilidad Presupuestal Sueldo Básico N° xxx del xx de xxxxx de 2020, Prima Vacaciones N° xxx del xx de xxxxx de 2020, Prima Servicio N° xxx del xx de xxxxx de 2020, Prima Navidad N° xxx del xx de xxxxx de 2020 y Cesantías N° xxx del xx de xxxxx de 2020.\n\nQue, en mérito de lo expuesto,\n\n"}
			articulo = Articulo{Texto: "Vincular para el TERCER PERIODO académico de 2020, en los programas de PREGRADO de la Universidad Distrital Francisco José de Caldas, en la modalidad de Hora Cátedra Prestaciones, en el escalafón y dedicación establecidas en la siguiente tabla, para el período comprendido entre el 28 de septiembre y hasta el 19 de diciembre del 2020, con recursos  asignados a los Certificados de Disponibilidad Presupuestal No Sueldo Básico N° xxx del xx de xxxxx de 2020, Prima Vacaciones N° xxx del xx de xxxxx de 2020, Prima Servicio N° xxx del xx de xxxxx de 2020, Prima Navidad N° xxx del xx de xxxxx de 2020 y Cesantías N° xxx del xx de xxxxx de 2020, a los siguientes docentes"}
		}
		if nivel == "POSGRADO" {
			resolucion = ResolucionCompleta{Consideracion: "Que mediante los artículos 3º y 4º del Decreto 1279 del 19 de Junio de 2002, el Gobierno Nacional estableció la naturaleza jurídica, laboral y prestacional de los docentes de las universidades estatales u oficiales del orden nacional, departamental, municipal y distrital, prescribiendo en su artículo 3º, que “[l]os profesores de hora-cátedra de las Universidades estatales u oficiales distintas a la Universidad Nacional de Colombia no son empleados públicos docentes de régimen especial ni pertenecen a la carrera profesoral y, por consiguiente, sus condiciones salariales y prestacionales no están regidas por el presente Decreto, sino por las reglas contractuales que en cada caso se convengan, conforme a las normas internas de cada Universidad, con sujeción a lo dispuesto en las disposiciones constitucionales y legales”.\n\nQue la Universidad Distrital Francisco José de Caldas puede vincular docentes de vinculación especial, en las modalidades de hora cátedra (HC), medio tiempo ocasional (MTO) y tiempo completo ocasional (TCO), Honorarios hora catedra, en virtud de lo consagrado en el artículo 13 del Acuerdo 011 de noviembre 15 de 2002 (Estatuto Docente), a término fijo, por periodos académicos.\n\nQue mediante Resolución Nro. 001 del 15 de febrero de 2012, proferida por la Vicerrectoría Académica, se establece el procedimiento para la selección y vinculación a la Universidad, de docentes de vinculación especial.\n\nQue los servicios de los docentes de vinculación especial, hora cátedra, de que trata el Acuerdo 011 de noviembre 15 de 2002, deberán ser reconocidos, en los términos del artículo 73 de la Ley 30 de 1992.\n\nQue el artículo 128 de la Constitución Política de Colombia, establece que nadie podrá desempeñar simultáneamente más de un empleo público, ni recibir más de una asignación que provenga del tesoro público, salvo las excepciones establecidas en la Ley.\n\nQue en virtud de la anterior norma constitucional, el literal d) del artículo 19 de la Ley 4ª de 1992 determinó, como excepción a dicha regla, entre otras, los honorarios percibidos por concepto de hora cátedra.\n\nQue la Corte Constitucional en Sentencia C-006 de 1996, a través de la cual se declararon inexequibles algunos apartes del artículo 73 de la Ley 30 de 1992, aclaró, entre otras cosas, que, la vinculación de los docentes de hora cátedra estará mediada por su transitoriedad o temporalidad, de lo cual se deriva la proporcionalidad prestacional al tiempo laborado.\n\nQue conforme al concepto de fecha diciembre 18 de 2018, emitido por el Departamento Administrativo de la Función Pública, “[e]n el caso de la Universidad Distrital Francisco José de Caldas, debe acogerse  a lo dispuesto en sus estatutos internos o en el acuerdo en el cual defina sus políticas y criterios generales para regular los procesos de selección, vinculación y contratación de los docentes ocasionales y catedráticos, en especial en lo relacionado con el reconocimiento y pago de prestaciones sociales y remuneración…”.\n\nQue conforme al parágrafo 1º del artículo quinto de la Resolución Nro. 001 de 2012, “[p]ara efectos de pago y liquidación el mes comprenderá (4) semanas o (30) días laborales en Proyectos Académicos de Pregrado y en los espacios académicos de Posgrado se pagará de acuerdo al número de horas efectivamente dictadas en cada período”.\n\nQue el Consejo Académico en uso de las funciones estatutarias establecidas en el Artículo 18, literales b.) y e.), del Estatuto General de la Universidad -Acuerdo 003 de abril 08 de 1997 expedido por el Consejo Superior Universitario-, mediante Resolución N° 015 de febrero 06 de 2020 aprobó el Calendario Académico del año 2020, para los programas académicos de posgrados de la Universidad Distrital Francisco José de Caldas.\n\nQue mediante resolución N°. 039 del 30 de junio del 2020 el Consejo Académico modifica el Artículo 1° de la Resolución 015 de febrero 06 de 2020 del calendario académico en programas de posgrado el cual quedo para iniciar clases hasta agosto 27 de 2020 y cierre periodo hasta el23 de diciembre de 2020.\n\nQue la Rectoría de la Universidad Distrital Francisco José de Caldas en uso de la autonomía universitaria y ante la problemática de salud pública, por medio de la Resolución 132 de marzo 19 de 2020, adoptó medidas transitorias y excepcionales de orden interno, de acuerdo con los lineamientos del orden Nacional y Distrital, para afrontar debidamente la situación de emergencia epidemiológica, evitando el contagio y garantizando la prestación del servicio a cargo de la Universidad.\n\nQue entre las medidas transitorias y excepcionales adoptadas en la Resolución de Rectoría 132 de marzo 19 de 2020, se optó por la suspensión de las clases presenciales en todos los programas pregrado y posgrado, a partir de las seis (6) de la mañana del 19 de marzo de 2020 y hasta nueva orden, según la actual coyuntura de salud pública, de emergencia económica y social.\n\nQue con base en la Directiva Ministerial del Ministerio de Educación Nacional Nº 04 del 22 de marzo de 2020, y con el propósito de dar continuidad durante el periodo de emergencia sanitaria a los programas académicos con registro calificado con que cuenta la Universidad Distrital Francisco José de Caldas en modalidad presencial, de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\nQue la misma Directiva Ministerial prevé que una vez finalizada la emergencia sanitaria, los desarrollos curriculares de los programas académicos con registro calificado en modalidad presencial deberán desarrollarse como se venía realizando habitualmente de acuerdo con las características propias de dicha modalidad.\n\nQue conforme lo dispuesto en el artículo 15 del Decreto Legislativo 491 de 2020, que tiene por objeto que las autoridades cumplan con la finalidad de proteger y garantizar los derechos y libertades de las personas, la primacía de los intereses generales, la sujeción de las autoridades a la Constitución y demás preceptos del ordenamiento jurídico, el cumplimiento de los fines y principios estatales, el funcionamiento eficiente y democrático de la administración y la observancia de los deberes del Estado y de los particulares, durante el período de aislamiento preventivo obligatorio, las universidades públicas dispondrán de las medidas necesarias para que los docentes ocasionales o de hora cátedra cumplan sus funciones mediante la modalidad de \"trabajo en casa\", entre otras, haciendo uso de las tecnologías de la información y las comunicaciones.\n\nQue se hace necesario garantizar la vinculación de docentes hora cátedra que realicen las actividades lectivas señaladas en el calendario académico y que no son suplidas por los docentes de carrera de la Universidad Distrital.\n\nQue en materia de Seguridad y Salud en el Trabajo SST para docentes ocasionales de la Universidad Distrital se deben informar las obligaciones específicas se encuentra contemplado en el Decreto 1072 de 2015 artículo 2.2.4.2.2.16. y la resolución de rectoría No. 624 de 2017, el docente ocasional, hora catedra y honorarios debe cumplir con las normas del sistema general de riesgos laborales.\n\nQue los servicios de los Docentes de Vinculación Especial señalados en el Acuerdo 011 de noviembre 15 de 2002 expedido por el Consejo Superior Universitario, deberán ser reconocidos mediante Resolución, (Resoluciones 0013 de enero 31 de 2003, 0013-A de enero 31 de 2003, Ley 30 de 1992 y Acuerdo 003 de 1997, Artículo 49 y ley 4 de 1992).\n\nQue, para efectos presupuestales, el presente acto administrativo se ejecutará con cargo a los recursos de que trata los Certificados de Disponibilidad Presupuestal No (s) Sueldo Básico N°, Prima Vacaciones N°, Prima Servicio N°, Prima Navidad N° Cesantías N°  XX del XX de enero de 2020.\n\n\n\nQue, en mérito de lo expuesto.\n\n"}
			articulo = Articulo{Texto: "VINCULAR: Vincular a los siguientes docentes para el TERCER PERIODO académico de 2020, en los programas de posgrado de la Universidad Distrital Francisco José de Caldas, como docentes en la modalidad de Hora Cátedra prestaciones , en el escalafón y dedicación establecidas en la siguiente tabla, para el período comprendido entre el XX de agosto y hasta 23 de diciembre del 2020, con recursos  asignados los Certificados de Disponibilidad Presupuestal No (s) Sueldo Básico N°, Prima Vacaciones N°, Prima Servicio N°, Prima Navidad N° Cesantías N°  XX del XX de enero de 2020."}
		}
		// accion = "Vincular"
		//nombreDedicacion = "Hora Cátedra"
	case "HCH":
		if nivel == "PREGRADO" {
			resolucion = ResolucionCompleta{Consideracion: "Que, con carácter general, el artículo 128 de la Constitución Política establece que “nadie podrá desempeñar simultáneamente más de un empleo público ni recibir más de una asignación que provenga del tesoro público, o de empresas o de instituciones en las que tenga parte mayoritaria el Estado, salvo los casos expresamente determinados por la ley (…) Entiéndese por tesoro público el de la Nación, el de las entidades territoriales y el de las descentralizadas”.\n\nQue, en desarrollo de este precepto superior, el artículo 19 de la Ley 4ª de 1992, en lo pertinente, establece que “nadie podrá desempeñar simultáneamente más de un empleo público, ni recibir más de una asignación que provenga del Tesoro Público, o de empresas o de instituciones en las que tenga parte mayoritaria el Estado. Exceptúense las siguientes asignaciones: (…) d) Los honorarios percibidos por concepto de hora-cátedra…”. En desarrollo de este acápite normativo, en particular, el parágrafo de la norma en cita prevé que “[n]o se podrán recibir honorarios que sumados correspondan a más de ocho (8) horas diarias de trabajo a varias entidades”.\n\nQue la Corte Constitucional, en Sentencia C-133 de 1993, señaló que, “si bien es cierto que en el artículo 128 C.P. se consagra una incompatibilidad, no lo es menos que ésta se encuentra en íntima relación de conexidad con la remuneración de los servidores estatales; basta ver que en ella se prohíbe la concurrencia de dos o más cargos públicos en una misma persona, tanto como recibir más de una asignación que provenga del erario público. El término ‘asignación´’ comprende toda clase de remuneración que emane del tesoro público, llámese sueldo, honorario, mesada pensional, etc. …”.\n\nQue, de conformidad con las normas y la jurisprudencia expuestas, la persona que tiene la calidad de pensionado del sector público podrá percibir otra asignación del Tesoro Público, siempre que la misma provenga de las excepciones establecidas en el artículo 19 de la Ley 4ª de 1992, como es el caso de los honorarios percibidos por los docentes que presten el servicio a una universidad o institución de educación del Estado, mediante el sistema de hora cátedra honorarios.\n\nQue mediante Resolución Nro. 001 del 15 de febrero de 2012, proferida por la Vicerrectoría Académica, se establece el procedimiento para la selección y vinculación a la Universidad, de docentes de vinculación especial.\n\nQue se hace necesario garantizar el reconocimiento de honorarios a docentes de hora cátedra en virtud de la Ley 4ª de 1992, que realicen actividades de horas lectivas señaladas en el calendario académico y que no son suplidas por los docentes de carrera de la Universidad Distrital.\n\nQue, mediante Acuerdo 002 de marzo 17 de 2011, emanado por el Consejo Superior Universitario, se modificó el artículo 43 del Acuerdo 011 de 2002, de manera que se pueden reconocer honorarios a personal pensionado en la modalidad de docentes de hora cátedra por honorarios.\n\nQue, al no existir norma expresa al respecto, los servicios de los docentes catedráticos que ostentan la condición de servidores públicos o pensionados, serán reconocidos mediante resolución.\n\nQue, para efectos del pago de los servicios prestados por los docentes hora catedra honorarios, se hará de conformidad con las horas lectivas efectivamente dictadas y aprobadas en los diferentes programas, así como en concordancia con el calendario académico emitido mediante Resolución 053 del veintitrés (23) de julio del 2020 para programas de pregrado.\n\nQue a los funcionarios públicos y pensionados del sector público se reconocerán honorarios hora catedra, de conformidad con lo dispuesto en la Resolución 01 de febrero del 2012 de la Vicerrectoría Académica y hasta por un máximo de (8) horas semanales en programas de pregrado.\n\nQue, dado que el vínculo que se genera con el docente catedrático no tiene carácter laboral, está obligado a realizar los aportes a seguridad social por su cuenta.\n\nQue el artículo 2 del Decreto 310 de 27 febrero 2020, emanado por el Presidente de la República de Colombia, establece, “ (…) a  partir del 1° de enero del 2020, fijase el valor del punto para empleados públicos docentes a quienes se les aplica el Decreto 1279 de 2002 y demás disposiciones que lo modifiquen o adicionen en catorce mil novecientos treinta y ocho pesos ($14.938) moneda corriente (…)”.\n\nQue, mediante Resolución 137 de abril 3 de 2020, emitida por Rectoría de la Universidad Distrital Francisco José de Caldas, se determinó “acoger el artículo 2 del Decreto 310 de 27 de febrero de 2020, en lo pertinente, única y expresamente al valor del punto para los docentes de Vinculación Especial hora cátedra, honorarios, Medio Tiempo Ocasional y Tiempo Completo Ocasional en pregrado, el valor del punto en CATORCE MIL NOVECIENTOS TREINTA Y OCHO PESOS ($14.938) MONEDA CORRIENTE”.\n\nQue, en materia de Seguridad y Salud en el Trabajo (SST), para docentes de hora cátedra honorarios de la Universidad Distrital Francisco José de Caldas, se deben informar las obligaciones específicas de que trata el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015, así como la Resolución de Rectoría No. 624 de 2017, conforme a los cuales, el docente ocasional, de hora catedra y por honorarios, debe cumplir con las normas del Sistema General de Riesgos Laborales. \n\nQue mediante los artículos 3º y 4º del Decreto 1279 del 19 de Junio de 2002, el Gobierno Nacional estableció la naturaleza jurídica, laboral y prestacional de los docentes de las universidades estatales u oficiales del orden nacional, departamental, municipal y distrital, prescribiendo en su artículo 3º, que “los profesores de hora-cátedra de las Universidades estatales u oficiales distintas a la Universidad Nacional de Colombia no son empleados públicos docentes de régimen especial ni pertenecen a la carrera profesoral y, por consiguiente, sus condiciones salariales y prestacionales no están regidas por el presente Decreto, sino por las reglas contractuales que en cada caso se convengan, conforme a las normas internas de cada Universidad, con sujeción a lo dispuesto en las disposiciones constitucionales y legales”.\n\nQue el Consejo Académico en uso de las funciones estatutarias establecidas en el Artículo 18, literales b.) y e.), del Estatuto General de la Universidad -Acuerdo 003 de abril 08 de 1997, expedido por el Consejo Superior Universitario-, mediante Resolución N°  053 de julio 23 de 2020 aprobó “el Calendario Académico especial para el segundo semestre del año 2020, para los programas académicos de pregrado bajo la modalidad virtual y/o alternancia, en el marco de las medidas de emergencia sanitaria y de aislamiento social decretado por el Gobierno Nacional, por cuenta de la pandemia Covid-19”\n\nQue la Resolución N°. 053 del 23 de julio del 2020, el Consejo Académico determinó en el artículo primero que las clases se inician el día 28 de septiembre de 2020 y se suspenden a partir del 20 de diciembre de 2020, cuando se entra en un periodo de receso hasta el 24 de enero de 2021.\n\nQue la Rectoría de la Universidad Distrital Francisco José de Caldas en uso de la autonomía universitaria y ante la problemática de salud pública, por medio de la Resolución 132 de marzo 19 de 2020, adoptó medidas transitorias y excepcionales de orden interno, de acuerdo con los lineamientos del orden Nacional y Distrital, para afrontar debidamente la situación de emergencia epidemiológica, evitando el contagio y garantizando la prestación del servicio a cargo de la Universidad.\n\nQue entre las medidas transitorias y excepcionales adoptadas en la Resolución de Rectoría 132 de marzo 19 de 2020, se optó por la “suspensión de las clases presenciales en todos los programas pregrado y posgrado de conformidad con el comunicado del 15 de marzo de 2020, a partir de las seis (6) de la mañana del 19 de marzo de 2020 y hasta nuevas instrucciones, según la actual coyuntura de salud pública, así como de emergencia económica y social”.\n\nQue con base en la Directiva Ministerial Nº 04 del 22 de marzo de 2020 del Ministerio de Educación Nacional, y con el propósito de dar continuidad durante el periodo de emergencia sanitaria a los programas académicos con registro calificado con que cuenta la Universidad Distrital Francisco José de Caldas en modalidad presencial, de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\nQue la misma Directiva Ministerial prevé que una vez finalizada la emergencia sanitaria, los desarrollos curriculares de los programas académicos con registro calificado en modalidad presencial deberán desarrollarse como se venía realizando habitualmente de acuerdo con las características propias de dicha modalidad.\n\nQue conforme lo dispuesto en el artículo 15 del Decreto Legislativo 491 de 2020, expedido por el Presidente de la República de Colombia, tiene por objeto que las autoridades cumplan con la finalidad de proteger y garantizar los derechos y libertades de las personas, la primacía de los intereses generales, la sujeción de las autoridades a la Constitución y demás preceptos del ordenamiento jurídico, el cumplimiento de los fines y principios estatales, el funcionamiento eficiente y democrático de la administración y la observancia de los deberes del Estado y de los particulares, durante el período de aislamiento preventivo obligatorio, las universidades públicas dispondrán de las medidas necesarias para que los docentes ocasionales o de hora cátedra cumplan sus funciones mediante la modalidad de \"trabajo en casa\", entre otras, haciendo uso de las tecnologías de la información y las comunicaciones.\n\n" +
				"Que el gasto que ocasione la presente resolución se hará con cargo a la disponibilidad presupuestal No. xxxx del xx de xxxxxx de 2020.\n\nQue, en mérito de lo expuesto, \n\n"}
			articulo = Articulo{Texto: "RECONOCIMIENTO: Reconocer Honorarios de conformidad con la Ley 4ª de 1992 para el tercer periodo académico del año 2020, en los programas de pregrado de la Universidad Distrital Francisco José de Caldas, como docentes en la modalidad de Hora Cátedra Honorarios (HCH), en el escalafón y dedicación establecidas en la siguiente tabla, para el período comprendido entre el XX de septiembre y hasta el 19 de diciembre del 2020, con recursos  asignados en el Certificado de Disponibilidad Presupuestal No. xxxx del xx de xxxxxx del 2020, a los siguientes docentes"}
		}
		if nivel == "POSGRADO" {
			resolucion = ResolucionCompleta{Consideracion: "Que mediante los artículos 3º y 4º del Decreto 1279 del 19 de Junio de 2002, el Gobierno Nacional estableció la naturaleza jurídica, laboral y prestacional de los docentes de las universidades estatales u oficiales del orden nacional, departamental, municipal y distrital, prescribiendo en su artículo 3º, que ”[l]os profesores de hora-cátedra de las Universidades estatales u oficiales distintas a la Universidad Nacional de Colombia no son empleados públicos docentes de régimen especial ni pertenecen a la carrera profesoral y, por consiguiente, sus condiciones salariales y prestacionales no están regidas por el presente Decreto, sino por las reglas contractuales que en cada caso se convengan, conforme a las normas internas de cada Universidad, con sujeción a lo dispuesto en las disposiciones constitucionales y legales”.\n\nQue el artículo 128 de la Carta Política, consigna que:\n\n”…nadie podrá desempeñar simultáneamente más de un empleo público ni recibir más de una asignación que provenga del tesoro público salvo las excepciones establecidas en la Ley…”\n\nQue dando alcance a la norma ibídem, la Ley 4 de 1992, en su artículo 19 determinó:\n\nNadie podrá desempeñar simultáneamente más de un empleo público, ni recibir más de una asignación que provenga del Tesoro Público, o de empresas o de instituciones en las que tenga parte mayoritaria el Estado. Exceptúense las siguientes asignaciones:\n\n(…)\n\nd.	Los honorarios percibidos por concepto de hora-cátedra\n\nEn desarrollo de este acápite normativo, en particular, el parágrafo de la norma en cita prevé que ”[n]o se podrán recibir honorarios que sumados correspondan a más de ocho (8) horas diarias de trabajo a varias entidades”.\n\nQue la Corte Constitucional, en Sentencia C-133 de 1993, señaló que, ”[s]i bien es cierto que en el artículo 128 C.P. se consagra una incompatibilidad, no lo es menos que ésta se encuentra en íntima relación de conexidad con la remuneración de los servidores estatales; basta ver que en ella se prohíbe la concurrencia de dos o más cargos públicos en una misma persona, tanto como recibir más de una asignación que provenga del erario público. El término ‘asignación´’ comprende toda clase de remuneración que emane del tesoro público, llámese sueldo, honorario, mesada pensional, etc. …”.\n\nQue, de conformidad con las normas y la jurisprudencia expuestas, la persona que tiene la calidad de pensionado del sector público podrá percibir otra asignación del Tesoro Público, siempre que la misma provenga de las excepciones establecidas en el artículo 19 de la Ley 4ª de 1992, como es el caso de los honorarios percibidos por los docentes que presten el servicio a una universidad o institución de educación del Estado, mediante el sistema de hora cátedra honorarios.\n\nQue mediante Resolución Nro. 001 del 15 de febrero de 2012, proferida por la Vicerrectoría Académica, se establece el procedimiento para la selección y vinculación a la Universidad, de docentes de vinculación especial.\n\nQue se hace necesario garantizar el reconocimiento de honorarios a docentes de hora cátedra en virtud de la Ley 4ª de 1992, que realicen actividades de horas lectivas señaladas en el calendario académico y que no son suplidas por los docentes de carrera de la Universidad Distrital.\n\nQue mediante el artículo 2° del Acuerdo 006 del 2002, que modificó el artículo 1° del Acuerdo 007 del 2001, ambos del Consejo Superior Universitario, se estableció que a los docentes hora catedra que pertenecen a la carrera docente de la Universidad Distrital Francisco José de Caldas, se les podrá reconocer hasta un máximo de seis (6) horas semanales adicionales a su carga normal, en posgrado.\n\nQue, el Acuerdo 005 del 27 de julio del 2001, mediante el cual fija el valor de la hora cátedra y se establece el número máximo de horas para los docentes que prestan servicios a la Universidad Distrital, en los programas de posgrado y su liquidación se efectúa con base al salario mínimo mensual legal vigente.\n\nQue los docentes desempeñen como docente de carrera de La Universidad Distrital deberán cumplir en el plan de trabajo con la carga académica de doce (12) horas semanales siendo aprobada por el Decano y coordinador.\n\nQue a los funcionarios público y pensionados, se les reconocerán honorarios hora catedra, de conformidad con lo dispuesto en la Resolución 01 de febrero del 2012 y hasta por un máximo de (8) horas semanales.\n\nQue, mediante Acuerdo 002 de marzo 17 de 2011, se modificó el artículo 43 del Acuerdo 011 de 2002, de manera que se pueden reconocer honorarios a personal pensionado en la modalidad de docentes de hora cátedra por honorarios.\n\nQue, dado que el presente reconocimiento no genera relación laboral con el docente, estos se encuentran obligados a realizar los aportes a seguridad social en forma independiente.\n\nQue el Consejo Académico en uso de las funciones estatutarias establecidas en el Artículo 18, literales b.) y e.), del Estatuto General de la Universidad -Acuerdo 003 de abril 08 de 1997 expedido por el Consejo Superior Universitario-, mediante Resolución N° 015 de febrero 06 de 2020 aprobó el Calendario Académico del año 2020, para los programas académicos de posgrados de la Universidad Distrital Francisco José de Caldas.\n\nQue mediante resolución N°. 039 del 30 de junio del 2020 el Consejo Académico modifica el Artículo 1° de la Resolución 015 de febrero 06 de 2020 del calendario académico en programas de posgrado el cual quedo para iniciar clases hasta agosto 27 de 2020 y cierre periodo hasta el23 de diciembre de 2020\n\nQue la Rectoría de la Universidad Distrital Francisco José de Caldas en uso de la autonomía universitaria y ante la problemática de salud pública, por medio de la Resolución 132 de marzo 19 de 2020, adoptó medidas transitorias y excepcionales de orden interno, de acuerdo con los lineamientos del orden Nacional y Distrital, para afrontar debidamente la situación de emergencia epidemiológica, evitando el contagio y garantizando la prestación del servicio a cargo de la Universidad.\n\nQue entre las medidas transitorias y excepcionales adoptadas en la Resolución de Rectoría 132 de marzo 19 de 2020, se optó por la suspensión de las clases presenciales en todos los programas pregrado y posgrado, a partir de las seis (6) de la mañana del 19 de marzo de 2020 y hasta nueva orden, según la actual coyuntura de salud pública, de emergencia económica y social.\n\nQue con base en la Directiva Ministerial del Ministerio de Educación Nacional Nº 04 del 22 de marzo de 2020, y con el propósito de dar continuidad durante el periodo de emergencia sanitaria a los programas académicos con registro calificado con que cuenta la Universidad Distrital Francisco José de Caldas en modalidad presencial, de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\nQue la misma Directiva Ministerial prevé que una vez finalizada la emergencia sanitaria, los desarrollos curriculares de los programas académicos con registro calificado en modalidad presencial deberán desarrollarse como se venía realizando habitualmente de acuerdo con las características propias de dicha modalidad.\n\nQue conforme lo dispuesto en el artículo 15 del Decreto Legislativo 491 de 2020, que tiene por objeto que las autoridades cumplan con la finalidad de proteger y garantizar los derechos y libertades de las personas, la primacía de los intereses generales, la sujeción de las autoridades a la Constitución y demás preceptos del ordenamiento jurídico, el cumplimiento de los fines y principios estatales, el funcionamiento eficiente y democrático de la administración y la observancia de los deberes del Estado y de los particulares, durante el período de aislamiento preventivo obligatorio, las universidades públicas dispondrán de las medidas necesarias para que los docentes ocasionales o de hora cátedra cumplan sus funciones mediante la modalidad de \"trabajo en casa\", entre otras, haciendo uso de las tecnologías de la información y las comunicaciones.\n\nQue el gasto que ocasione la presente resolución se hará con cargo a la disponibilidad presupuestal No. XXXX del XX DE XXXXXXX del 2020.\n\nQue, en mérito de lo expuesto,\n\n"}
			articulo = Articulo{Texto: "RECONOCIMIENTO: Reconocer Honorarios de conformidad con la Ley 4ª de 1992 para el tercer periodo académico del año 2020, en los programas de posgrado de la Universidad Distrital Francisco José de Caldas, como docentes en la modalidad de Hora Cátedra Honorarios (HCH), en el escalafón y dedicación establecidas en la siguiente tabla, para el período comprendido entre el XX de agosto y hasta 23 de diciembre del 2020, con recursos  asignados en el Certificado de Disponibilidad Presupuestal No (s)  XX del XX de enero de 2020"}
		}
		// accion = "Reconocer Honorarios"
		//nombreDedicacion = "Hora Cátedra Honorarios"
	case "TCO-MTO":
		resolucion = ResolucionCompleta{Consideracion: "Que el artículo tercero del Decreto 1279 de junio 19 de 2002, mediante el cual se establece el régimen salarial y prestacional de los docentes de las universidades estatales, señala que “[l]os profesores ocasionales no son empleados públicos docentes de régimen especial ni pertenecen a la carrera profesoral y, por consiguiente, sus condiciones salariales y prestacionales no están regidas por el presente Decreto”, precisando que, “[n]o obstante, su vinculación se hace conforme a las reglas que define cada Universidad, con sujeción a lo dispuesto por la ley 30 de 1992 y demás disposiciones constitucionales y legales vigentes”.\n\nQue la Universidad Distrital Francisco José de Caldas puede vincular docentes de vinculación especial en las modalidades de hora cátedra (HC), medio tiempo ocasional (MTO) y tiempo completo ocasional (TCO), de que trata el artículo 13 del Acuerdo 011 de noviembre 15 de 2002 (Estatuto Docente), a término fijo, por periodos académicos.\n\nQue, mediante Resolución Nro. 001 del 15 de febrero de 2012, proferida por la Vicerrectoría Académica, se estableció el procedimiento para la selección y vinculación, a la Universidad, de docentes de vinculación especial.\n\nQue los servicios de los docentes de vinculación especial, deberán ser reconocidos en los términos del inciso segundo del artículo 74 de la Ley 30 de 1992, esto es, “mediante resolución”.\n\nQue la Corte Constitucional, en Sentencia C-006 de 1996, a través de la cual se declararon inexequibles algunos apartes del artículo 73 de la Ley 30 de 1992, aclaró, entre otras cosas, que la vinculación de los docentes de vinculación especial estará mediada por su transitoriedad o temporalidad, de lo cual se deriva la proporcionalidad prestacional al tiempo laborado.\n\nQue, conforme al concepto de fecha diciembre 18 de 2018, emitido por el Departamento Administrativo de la Función Pública, “[e]n el caso de la Universidad Distrital Francisco José de Caldas, debe acogerse  a lo dispuesto en sus estatutos internos o en el acuerdo en el cual defina sus políticas y criterios generales para regular los procesos de selección, vinculación y contratación de los docentes ocasionales y catedráticos, en especial en lo relacionado con el reconocimiento y pago de prestaciones sociales y remuneración…”.\n\nQue el parágrafo 1º del artículo quinto de la Resolución Nro. 001 de 2012 de la Vicerrectoría Académica, establece que,“[p]ara efectos de pago y liquidación el mes comprenderá (4) semanas o (30) días laborales en Proyectos Académicos de Pregrado y en los espacios académicos de Posgrado se pagará de acuerdo al número de horas efectivamente dictadas en cada período”.\n\nQue el artículo 2 del Decreto 310 de 27 febrero 2020, emanado por el Presidente de la República de Colombia, establece, “ (…) a  partir del 1° de enero del 2020, fijase el valor del punto para empleados públicos docentes a quienes se les aplica el Decreto 1279 de 2002 y demás disposiciones que lo modifiquen o adicionen en catorce mil novecientos treinta y ocho pesos ($14.938) moneda corriente (…)”.\n\nQue, mediante Resolución 137 de abril 3 de 2020, emitida por Rectoría de la Universidad Distrital Francisco José de Caldas, se determinó “acoger el artículo 2 del Decreto 310 de 27 de febrero de 2020, en lo pertinente, única y expresamente al valor del punto para los docentes de Vinculación Especial hora cátedra, honorarios, Medio Tiempo Ocasional y Tiempo Completo Ocasional en pregrado, el valor del punto en CATORCE MIL NOVECIENTOS TREINTA Y OCHO PESOS ($14.938) MONEDA CORRIENTE”.\n\nQue, en materia de Seguridad y Salud en el Trabajo (SST), para docentes ocasionales de la Universidad Distrital Francisco José de Caldas, se deben informar las obligaciones específicas de que trata el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015, así como la Resolución de Rectoría No. 624 de 2017, conforme a los cuales, el docente ocasional, de hora catedra y por honorarios, debe cumplir con las normas del Sistema General de Riesgos Laborales.\n\nQue el Consejo Académico en uso de las funciones estatutarias establecidas en el Artículo 18, literales b.) y e.), del Estatuto General de la Universidad -Acuerdo 003 de abril 08 de 1997, expedido por el Consejo Superior Universitario-, mediante Resolución N°  053 de julio 23 del 2020 aprobó “el Calendario Académico especial para el segundo semestre del año 2020, para los programas académicos de pregrado bajo la modalidad virtual y/o alternancia, en el marco de las medidas de emergencia sanitaria y de aislamiento social decretado por el Gobierno Nacional, por cuenta de la pandemia Covid-19”\n\nQue la Resolución N°. 053 del 23 de julio del 2020, el Consejo Académico determinó en el artículo primero que las clases se inician el día 28 de septiembre de 2020 y se suspenden a partir del 20 de diciembre de 2020, cuando se entra en un periodo de receso hasta el 24 de enero de 2021.\n\nQue la Rectoría de la Universidad Distrital Francisco José de Caldas en uso de la autonomía universitaria y ante la problemática de salud pública, por medio de la Resolución 132 de marzo 19 de 2020, adoptó medidas transitorias y excepcionales de orden interno, de acuerdo con los lineamientos del orden Nacional y Distrital, para afrontar debidamente la situación de emergencia epidemiológica, evitando el contagio y garantizando la prestación del servicio a cargo de la Universidad.\n\nQue entre las medidas transitorias y excepcionales adoptadas en la Resolución de Rectoría 132 de marzo 19 de 2020, se optó por la “suspensión de las clases presenciales en todos los programas pregrado y posgrado de conformidad con el comunicado del 15 de marzo de 2020, a partir de las seis (6) de la mañana del 19 de marzo de 2020 y hasta nuevas instrucciones, según la actual coyuntura de salud pública, así como de emergencia económica y social”.\n\nQue con base en la Directiva Ministerial Nº 04 del 22 de marzo de 2020 del Ministerio de Educación Nacional, y con el propósito de dar continuidad durante el periodo de emergencia sanitaria a los programas académicos con registro calificado con que cuenta la Universidad Distrital Francisco José de Caldas en modalidad presencial, de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\nQue la misma Directiva Ministerial prevé que una vez finalizada la emergencia sanitaria, los desarrollos curriculares de los programas académicos con registro calificado en modalidad presencial deberán desarrollarse como se venía realizando habitualmente de acuerdo con las características propias de dicha modalidad.\n\nQue conforme lo dispuesto en el artículo 15 del Decreto Legislativo 491 de 2020, expedido por el Presidente de la República de Colombia, tiene por objeto que las autoridades cumplan con la finalidad de proteger y garantizar los derechos y libertades de las personas, la primacía de los intereses generales, la sujeción de las autoridades a la Constitución y demás preceptos del ordenamiento jurídico, el cumplimiento de los fines y principios estatales, el funcionamiento eficiente y democrático de la administración y la observancia de los deberes del Estado y de los particulares, durante el período de aislamiento preventivo obligatorio, las universidades públicas dispondrán de las medidas necesarias para que los docentes ocasionales o de hora cátedra cumplan sus funciones mediante la modalidad de "trabajo en casa", entre otras, haciendo uso de las tecnologías de la información y las comunicaciones.\n\nQue se hace necesario garantizar la vinculación de docentes hora cátedra que realicen las actividades lectivas señaladas en el calendario académico y que no son suplidas por los docentes de carrera de la Universidad Distrital\n\nQue, para efectos presupuestales, el presente acto administrativo se ejecutará con cargo a los recursos de que trata los Certificados de Disponibilidad Presupuestal No (s) Sueldo Básico N° XXX del xxxx de xxxxx de 2020, Prima Vacaciones N° XXX del xxxx de xxxxx de 2020, Prima Servicio N° XXX del xxxx de xxxxx de 2020, Prima Navidad N° XXX del xxxx de xxxxx de 2020, Cesantías N°  XXX del xxxx de xxxxx de 2020.\n\nEn mérito de lo expuesto,\n\n"}
		articulo = Articulo{Texto: "Vincular para el TERCER PERIODO académico de 2020 dentro periodo comprendido entre el veintiocho (28) de septiembre al diecinueve (19) de diciembre del 2020, en los programas de pregrado de la Universidad Distrital Francisco José de Caldas, en la modalidad de Medio Tiempo y Tiempo Completo Ocasional  en el escalafón y dedicación establecidas en la siguiente tabla, con los recursos  asignados en los Certificados de Disponibilidad Presupuestal No (s) Sueldo Básico N° XXX del xxxx de xxxxx de 2020, Prima Vacaciones N° XXX del xxxx de xxxxx de 2020, Prima Servicio N° XXX del xxxx de xxxxx de 2020, Prima Navidad N° XXX del xxxx de xxxxx de 2020, Cesantías N°  XXX del xxxx de xxxxx de 2020, a los siguientes docentes"}
		// accion = "Vincular"
		//nombreDedicacion = "Medio Tiempo Ocasional y Tiempo Completo Ocasional"
	}

	// if tipo != "1" {
	// 	articulo = Articulo{Texto: "Modificar la Resolución No XXX del XXX del XXXX en cuanto al número de horas semanales y el valor total para el " + periodoStr + " Período Académico del " + strconv.Itoa(vigencia) + ", como docentes en la modalidad de " + nombreDedicacion + " de Vinculación Especial, en el escalafón y dedicación establecidas en la siguiente tabla:"}
	// }

	articulos = append(articulos, articulo)

	/////ARTICULO 2
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "PREPARACIÓN DE CURRICULOS. El proceso de revisión, preparar de las clases y adaptación de los syllabus de cada espacio académico con la apropiación de herramientas digitales como instrumentos de apoyo y demás actividades se harán con base en la distribución horaria establecida por cada Facultad.\n\n"}
		paragrafo := Paragrafo{Texto: "De conformidad con los considerados resolución N°. 039 del 30 de junio del 2020, mediante el la cual se fija el calendario académico para el tercer periodo académico de la actual anualidad, los docentes podrán apoyarse en Planestic y demás instancias funcionales de la Universidad, para llevar a cabo el proceso de revisión y adaptación de las asignaturas para el uso de herramientas digitales de apoyo, y la preparación de los diferentes contenidos temáticos."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "REMUNERACIÓN. El pago de los servicios prestados por los profesores de vinculación especial a que se refiere el artículo anterior, según su escalafón, se hará previa certificación de las horas efectivamente dictadas, que se encuentren inmersas en el correspondiente plan de trabajo de la gestión académica, expedida por el decano y/o director de proyecto curricular.\n\n"}
		paragrafo := Paragrafo{Texto: "El valor del salario mínimo mensual legal vigente para el reconocimiento y pago de los docentes en cuestión, será el que fije el Gobierno Nacional."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "TCO-MTO" {
		//MTO Y TCO
		articulo = Articulo{Texto: "El pago de los servicios prestados por los profesores de vinculación especial a que se refiere el artículo anterior, según su escalafón, se hará previa certificación de las horas efectivamente dictadas, que se encuentren inmersas en el correspondiente plan de trabajo de la gestión académica, expedida por el decano y/o director de proyecto curricular.\n\n"}
		paragrafo := Paragrafo{Texto: "El valor del punto en pesos para el reconocimiento y pago de los docentes en cuestión, será el que fije el Gobierno Nacional, mediante decreto, cada año, y que la Universidad Distrital Francisco José de Caldas acoja, mediante acto administrativo."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "HCH" && nivel == "PREGRADO" {
		articulo = Articulo{Texto: "El pago de los honorarios por los servicios prestados a los catedráticos a que se refiere el artículo anterior, según su escalafón, se hará previa certificación de las horas efectivamente dictadas, que se encuentren inmersas en el correspondiente plan de trabajo de la gestión académica, expedida por el decano y/o director de proyecto curricular.\n\n"}
		paragrafo := Paragrafo{Texto: "El valor del punto en pesos para el reconocimiento y pago de los docentes de hora cátedra honorarios, será el que fije el Gobierno Nacional mediante decreto, cada año, y que la Universidad Distrital Francisco José de Caldas acoja mediante acto administrativo."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "HCP" && nivel == "PREGRADO" {
		articulo = Articulo{Texto: "El pago de los servicios prestados por los profesores de vinculación especial a que se refiere el artículo anterior, según su escalafón, se hará previa certificación de las horas efectivamente dictadas, que se encuentren inmersas en el correspondiente plan de trabajo de la gestión académica, expedida por el decano y/o director de proyecto curricular.\n\n"}
		paragrafo := Paragrafo{Texto: "El valor del punto en pesos para el reconocimiento y pago de los docentes en cuestión, será el que fije el Gobierno Nacional mediante decreto cada año y que la Universidad Distrital Francisco José de Caldas acogerá mediante acto administrativo."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}

	articulos = append(articulos, articulo)
	/////ARTICULO 3
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "OBLIGACIONES. Los docentes en cuestión deberán cumplir con las obligaciones inherentes a la naturaleza del servicio, contempladas en la ley, en los reglamentos de la Universidad Distrital y en los planes de trabajo establecidos en el aplicativo de gestión académica, entregados por cada profesor, y aprobados por el coordinador del correspondiente programa académico, decano y/o director.\n\n"}
		paragrafo = Paragrafo{Texto: " Las funciones propias de la vinculación docente serán aquellas que se fijen en el Plan de Trabajo del período académico correspondiente y de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\n"}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		paragrafo = Paragrafo{Texto: "El docente cumplirá con lo establecido en la Circular No. 046 de julio 19 de 2017 de Rectoría, con relación al pago de aportes a seguridad social de forma independiente, de conformidad con la ley."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "PREPARACIÓN DE CURRICULOS. El proceso de revisión, preparar de las clases y adaptación de los syllabus de cada espacio académico con la apropiación de herramientas digitales como instrumentos de apoyo y demás actividades se harán con base en la distribución horaria establecida por cada Facultad.\n\n"}
		paragrafo := Paragrafo{Texto: "De conformidad con los considerados resolución N°. 039 del 30 de junio del 2020, mediante el la cual se fija el calendario académico para el tercer periodo académico de la actual anualidad, los docentes podrán apoyarse en Planestic y demás instancias funcionales de la Universidad, para llevar a cabo el proceso de revisión y adaptación de las asignaturas para el uso de herramientas digitales de apoyo, y la preparación de los diferentes contenidos temáticos."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "TCO-MTO" {
		//MTO Y TCO
		articulo = Articulo{Texto: "Los docentes en cuestión deberán cumplir con las obligaciones inherentes a la naturaleza del servicio, contempladas en la ley, en los reglamentos de la Universidad Distrital Francisco José de Caldas y en los planes de trabajo establecidos en el aplicativo de gestión académica, entregados por cada profesor, y aprobados por el coordinador del correspondiente programa académico, decano y/o director.\n\n"}
		paragrafo = Paragrafo{Texto: "Las funciones propias de la vinculación docente serán aquellas que se fijen en el respectivo plan de trabajo."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}
	if dedicacion == "HCP" && nivel == "PREGRADO" {
		articulo = Articulo{Texto: "Los docentes en cuestión deberán cumplir con las obligaciones inherentes a la naturaleza del servicio, contempladas en la ley, en los reglamentos de la Universidad Distrital y en los planes de trabajo establecidos en el aplicativo de gestión académica, entregados por cada profesor, y aprobados por el coordinador del correspondiente programa académico, decano y/o director\n\n"}
		paragrafo = Paragrafo{Texto: "Las funciones propias de la vinculación docente serán aquellas que se fijen en el Plan de Trabajo del período académico correspondiente y de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	}

	if dedicacion == "HCH" && nivel == "PREGRADO" {

		// Pregrado y HCH
		articulo = Articulo{Texto: "OBLIGACIONES. Los docentes en cuestión deberán cumplir con las obligaciones inherentes a la naturaleza del servicio, contempladas en la ley, en los reglamentos de la Universidad Distrital y en los planes de trabajo establecidos en el aplicativo de gestión académica, entregados por cada profesor, y aprobados por el coordinador del correspondiente programa académico, decano y/o director.\n\n"}
		paragrafo = Paragrafo{Texto: "Las funciones propias de la vinculación docente serán aquellas que se fijen en el Plan de Trabajo del período académico correspondiente y de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad.\n\n"}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		paragrafo = Paragrafo{Texto: "El docente cumplirá con lo establecido en la Circular No. 046 de julio 19 de 2017 de Rectoría, con relación al pago de aportes a seguridad social de forma independiente, de conformidad con la ley."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)

	}
	articulos = append(articulos, articulo)
	////ARTICULO 4
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "TERMINACION. En caso de incumplimiento o retiro del docente, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, declarará la terminación del reconocimiento con corte a la fecha del respectivo cumplido, cancelando los correspondientes Honorarios, conforme al cálculo que efectúe la División de Recursos Humanos.\n\n"}
		articulos = append(articulos, articulo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "OBLIGACIONES. Los docentes en cuestión deberán cumplir con las obligaciones inherentes a la naturaleza del servicio, contempladas en la ley, en los reglamentos de la Universidad Distrital y en los planes de trabajo establecidos en el aplicativo de gestión académica, entregados por cada profesor, y aprobados por el coordinador del correspondiente programa académico, decano y/o director.\n\n"}
		paragrafo := Paragrafo{Texto: "Las funciones propias de la vinculación docente serán aquellas que se fijen en el Plan de Trabajo del período académico correspondiente y de manera excepcional se podrá desarrollar actividades académicas asistidas por las herramientas que ofrecen las Tecnologías de la Información y las Comunicaciones TIC, garantizando las condiciones de calidad reconocidas en el registro calificado, sin que esto implique el cambio de modalidad."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		paragrafo = Paragrafo{Texto: "El docente cumplirá con lo establecido en la Circular No. 046 de julio 19 de 2017 de Rectoría, con relación al pago de aportes a seguridad social de forma independiente, de conformidad con la ley."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		articulos = append(articulos, articulo)

	} else {
		if dedicacion == "HCP" && nivel == "PREGRADO" {
			articulo = Articulo{Texto: "En caso de incumplimiento o retiro del docente, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, declarará la terminación del vínculo y la liquidación con corte a la fecha del respectivo cumplido, cancelando las correspondientes prestaciones sociales, conforme al cálculo que efectúe la División de Recursos Humanos."}
			articulos = append(articulos, articulo)
		}
		if dedicacion == "HCH" && nivel == "PREGRADO" {
			// Pregrado y HCH
			articulo = Articulo{Texto: "TERMINACION. En caso de incumplimiento o retiro del docente, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, declarará la terminación del reconocimiento con corte a la fecha del respectivo cumplido, cancelando los correspondientes Honorarios, conforme al cálculo que efectúe la División de Recursos Humanos."}
			articulos = append(articulos, articulo)
		}
		if dedicacion == "TCO-MTO" {
			articulo = Articulo{Texto: "En caso de incumplimiento o retiro del docente, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, declarará la terminación del vínculo laboral y la liquidación con corte a la fecha del respectivo cumplido, cancelando las correspondientes prestaciones sociales, conforme al cálculo que efectúe la División de Recursos Humanos.\n\n"}
			articulos = append(articulos, articulo)
		}
	}
	/////ARTICULO 5
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "RECURSOS. - El gasto que ocasione el presente acto administrativo, se hará con cargo al presupuesto de la actual vigencia, previa certificación de disponibilidad presupuestal.\n\n"}
		paragrafo = Paragrafo{Texto: "En todo caso, los pagos correspondientes estarán sujetos a las apropiaciones presupuestales y a las transferencias realizadas por la Secretaría de Hacienda Distrital y se realizarán dentro de los primeros cinco (5) días hábiles de cada mes."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
	} else {
		if dedicacion == "HCH" && nivel == "PREGRADO" {
			// Pregrado y HCH
			articulo = Articulo{Texto: "RECURSOS. - El gasto que ocasione el presente acto administrativo, se hará con cargo al presupuesto de la actual vigencia, previa certificación de disponibilidad presupuestal.\n\n"}
			paragrafo = Paragrafo{Texto: "En todo caso, los pagos correspondientes estarán sujetos a las apropiaciones presupuestales y a las transferencias realizadas por la Secretaría de Hacienda Distrital y se realizarán dentro de los primeros cinco (5) días hábiles de cada mes."}
			articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		}
		if dedicacion == "HCP" && nivel == "POSGRADO" {
			articulo = Articulo{Texto: "TERMINACION. En caso de incumplimiento o retiro del docente, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, declarará la terminación del vínculo y la liquidación con corte a la fecha del respectivo cumplido, cancelando las correspondientes prestaciones sociales, conforme al cálculo que efectúe la División de Recursos Humanos.\n\n"}
		}
		if dedicacion == "HCP" && nivel == "PREGRADO" {
			//HCP Pregrado
			articulo = Articulo{Texto: "El gasto que ocasione el presente acto administrativo, se hará con cargo al presupuesto de la actual vigencia, previa certificación de disponibilidad presupuestal.\n\n"}
			paragrafo = Paragrafo{Texto: "En todo caso, los pagos correspondientes estarán sujetos a las apropiaciones presupuestales y a las transferencias realizadas por la Secretaría de Hacienda Distrital y se realizarán dentro de los primeros cinco (5) días hábiles de cada mes."}
			articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		}
		if dedicacion == "TCO-MTO" {
			//tco-mto
			articulo = Articulo{Texto: "El gasto que ocasione el presente acto administrativo, se hará con cargo al presupuesto de la actual vigencia, previa certificación de disponibilidad presupuestal.\n\n"}
			paragrafo = Paragrafo{Texto: "En todo caso, los pagos correspondientes estarán sujetos a las apropiaciones presupuestales y a las transferencias realizadas por la Secretaría de Hacienda Distrital, y se realizarán dentro de los primeros cinco (5) días hábiles de cada mes."}
			articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		}
	}
	articulos = append(articulos, articulo)
	/////ARTICULO 6
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "SUSPENSIÓN. -  En el supuesto de que se declare la suspensión de actividades académicas, por parte de los órganos competentes de la Universidad Distrital Francisco José de Caldas, cesará automáticamente para el docente de vinculación especial, la obligación de prestar sus servicios y para la Universidad, la de pagar los honorarios correspondientes al periodo suspendido. En este caso, la fecha de terminación de la vinculación se prolongará por un período igual al de la suspensión.\n\nSi la suspensión de actividades académicas persiste por más de quince (15) días calendario, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, podrá declarar la terminación del vínculo y efectuar la liquidación correspondiente, con corte a la fecha del respectivo cumplido, conforme al cálculo que efectúe la División de Recursos Humanos."}
		articulos = append(articulos, articulo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "RECURSOS. - El gasto que ocasione el presente acto administrativo, se hará con cargo al presupuesto de la actual vigencia, previa certificación de disponibilidad presupuestal.\n\n"}
		paragrafo = Paragrafo{Texto: "En todo caso, los pagos correspondientes estarán sujetos a las apropiaciones presupuestales y a las transferencias realizadas por la Secretaría de Hacienda Distrital y se realizarán dentro de los primeros cinco (5) días hábiles de cada mes."}
		articulo.Paragrafos = append(articulo.Paragrafos, paragrafo)
		articulos = append(articulos, articulo)
	} else {
		if dedicacion == "HCH" && nivel == "PREGRADO" {
			// Pregrado y HCH
			articulo = Articulo{Texto: "SUSPENSIÓN. -  En el supuesto de que se declare la suspensión de actividades académicas, por parte de los órganos competentes de la Universidad Distrital Francisco José de Caldas, cesará automáticamente para el docente de vinculación especial, la obligación de prestar sus servicios y para la Universidad, la de pagar los honorarios correspondientes al periodo suspendido. En este caso, la fecha de terminación de la vinculación se prolongará por un período igual al de la suspensión.\n\nSi la suspensión de actividades académicas persiste por más de quince (15) días calendario, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, podrá declarar la terminación del vínculo y efectuar la liquidación correspondiente, con corte a la fecha del respectivo cumplido, conforme al cálculo que efectúe la División de Recursos Humanos."}
			articulos = append(articulos, articulo)
		}
		if dedicacion == "HCP" && nivel == "PREGRADO" {
			//HCP
			articulo = Articulo{Texto: "En el supuesto de que se declare la suspensión de actividades académicas, por parte de los órganos competentes de la Universidad Distrital Francisco José de Caldas, cesará automáticamente para el docente de vinculación especial, la obligación de prestar sus servicios y para la Universidad, la de pagar los salarios y prestaciones correspondientes al periodo suspendido, pero persistirá, para esta última, la de efectuar los respectivos aportes a salud y pensión, en el porcentaje que le corresponda. En este caso, la fecha de terminación de la vinculación se prolongará por un período igual al de la suspensión\n\nSi la suspensión de actividades académicas persiste por más de quince (15) días calendario, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, podrá declarar la terminación del vínculo y efectuar la liquidación correspondiente, con corte a la fecha del respectivo cumplido, cancelando las correspondientes prestaciones sociales, conforme al cálculo que efectúe la División de Recursos Humanos"}
			articulos = append(articulos, articulo)
		}
		if dedicacion == "TCO-MTO" {
			//TCO-MTO
			articulo = Articulo{Texto: " En el supuesto de que se declare la suspensión de actividades académicas, por parte de los órganos competentes de la Universidad Distrital Francisco José de Caldas, cesará automáticamente para el docente de vinculación especial, la obligación de prestar sus servicios y para la Universidad, la de pagar los salarios y prestaciones correspondientes al periodo suspendido, pero persistirá, para esta última, la de efectuar los respectivos aportes a salud y pensión, en el porcentaje que le corresponda. En este caso, la fecha de terminación de la vinculación se prolongará por un período igual al de la suspensión.\n\nSi la suspensión de actividades académicas persiste por más de quince (15) días calendario, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, podrá declarar la terminación del vínculo y efectuar la liquidación correspondiente, con corte a la fecha del respectivo cumplido, cancelando las correspondientes prestaciones sociales, conforme al cálculo que efectúe la División de Recursos Humanos."}
			articulos = append(articulos, articulo)
		}
	}
	/////ARTICULO 7
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "RIESGOS LABORALES. -  Los docentes hora cátedra honorarios a que se refiere el presente acto administrativo, en los términos previstos en el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015 y la Resolución de Rectoría No. 624 de 2017, darán cumplimiento a las normas del Sistema General de Riesgos Laborales."}
		articulos = append(articulos, articulo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "SUSPENSIÓN. -  En el supuesto de que se declare la suspensión de actividades académicas, por parte de los órganos competentes de la Universidad Distrital Francisco José de Caldas, cesará automáticamente para el docente de vinculación especial, la obligación de prestar sus servicios y para la Universidad, la de pagar los honorarios correspondientes al periodo suspendido. En este caso, la fecha de terminación de la vinculación se prolongará por un período igual al de la suspensión.\n\nSi la suspensión de actividades académicas persiste por más de quince (15) días calendario, la Universidad Distrital Francisco José de Caldas, mediante acto administrativo motivado, podrá declarar la terminación del vínculo y efectuar la liquidación correspondiente, con corte a la fecha del respectivo cumplido, conforme al cálculo que efectúe la División de Recursos Humanos."}
		articulos = append(articulos, articulo)
	} else {
		if dedicacion == "TCO-MTO" {
			articulo = Articulo{Texto: "Los docentes ocasionales a que se refiere el presente acto administrativo, en los términos previstos en el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015 y la Resolución de Rectoría No. 624 de 2017, darán cumplimiento a las normas del Sistema General de Riesgos Laborales.\n\n"}
			articulos = append(articulos, articulo)
		} else {
			if dedicacion == "HCH" && nivel == "PREGRADO" {
				articulo = Articulo{Texto: "RIESGOS LABORALES. -  Los docentes hora cátedra honorarios a que se refiere el presente acto administrativo, en los términos previstos en el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015 y la Resolución de Rectoría No. 624 de 2017, darán cumplimiento a las normas del Sistema General de Riesgos Laborales."}
				articulos = append(articulos, articulo)
			}
			if dedicacion == "HCP" && nivel == "PREGRADO" {
				// Pregrado y HCP
				articulo = Articulo{Texto: "Los docentes hora cátedra a que se refiere el presente acto administrativo, en los términos previstos en el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015 y la Resolución de Rectoría No. 624 de 2017, darán cumplimiento a las normas del Sistema General de Riesgos Laborales."}
				articulos = append(articulos, articulo)
			}
		}
	}
	/////ARTICULO 8
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "INHABILIDAD O INCOMPATIBILIDAD.- Comuníquese la presente resolución a los docentes mencionados en el artículo primero, quienes deberán manifestar, bajo la gravedad de juramento que se entiende prestado con la aprobación del formato de carga en el aplicativo de gestión académica, que no se encuentran incursos en causal de inhabilidad o incompatibilidad establecida en las normas pertinentes y aplicables, que no tienen cruces de horarios y que el monto de los honorarios que perciben por concepto de hora cátedra no corresponden a más de ocho (8) horas diarias de trabajo a varias entidades."}
		articulos = append(articulos, articulo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "RIESGOS LABORALES. -  Los docentes hora cátedra a que se refiere el presente acto administrativo, en los términos previstos en el artículo 2.2.4.2.2.16. del Decreto 1072 de 2015 y la Resolución de Rectoría No. 624 de 2017, darán cumplimiento a las normas del Sistema General de Riesgos Laborales."}
		articulos = append(articulos, articulo)
	} else {
		if dedicacion == "HCH" && nivel == "PREGRADO" {
			// Pregrado y HCH
			articulo = Articulo{Texto: "INHABILIDAD O INCOMPATIBILIDAD.- Comuníquese la presente resolución a los docentes mencionados en el artículo primero, quienes deberán manifestar, bajo la gravedad de juramento que se entiende prestado con la aprobación del formato de carga en el aplicativo de gestión académica, que no se encuentran incursos en causal de inhabilidad o incompatibilidad establecida en las normas pertinentes y aplicables, que no tienen cruces de horarios y que el monto de los honorarios que perciben por concepto de hora cátedra no corresponden a más de ocho (8) horas diarias de trabajo a varias entidades."}
			articulos = append(articulos, articulo)
		} else {
			if dedicacion == "TCO-MTO" {
				articulo = Articulo{Texto: "Comuníquese la presente resolución a los docentes mencionados en el artículo primero, quienes deberán manifestar, bajo la gravedad de juramento que se entiende prestado con la aprobación del formato de carga en el aplicativo de gestión académica, que no se encuentran incursos en causal de inhabilidad o incompatibilidad, establecida en las normas pertinentes y aplicables, así como que no tienen cruces de horarios."}
				articulos = append(articulos, articulo)
			}
			if dedicacion == "HCP" && nivel == "PREGRADO" {
				//HCPPregrado y Posgrado
				articulo = Articulo{Texto: "Comuníquese la presente resolución a los docentes mencionados en el artículo primero, quienes deberán manifestar, bajo la gravedad de juramento que se entiende prestado con la aprobación del formato de carga en el aplicativo de gestión académica, que no se encuentran incursos en causal de inhabilidad o incompatibilidad establecida en las normas pertinentes y aplicables, que no tienen cruces de horarios."}
				articulos = append(articulos, articulo)
			}
		}
	}
	/////ARTICULO 9
	// Honorarios
	if dedicacion == "HCH" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "VIGENCIA. -  El presente acto administrativo se expide a los XXX (XXX) días del mes de XXX del año dos mil XXX (20XX) y surte efectos de conformidad con la Resolución No. 039 de junio treinta (30) de 2020, por medio del cual se expide el calendario académico para el año 2020, en concreto, para las actividades académicas relacionadas con el tercer periodo académico del año 2020, que corresponde del xxxxxx (00) de agosto y hasta 23 de diciembre del 2020."}
		articulos = append(articulos, articulo)
	}
	if dedicacion == "HCH" && nivel == "PREGRADO" {
		articulo = Articulo{Texto: "VIGENCIA. -  El presente acto administrativo se expide a los XXX (XXX) días del mes de XXX del año dos mil XXX (20XX) y surte efectos de conformidad con la Resolución No. 053 de julio veintitrés (23) de 2020, por medio del cual se expide el calendario académico para el segundo semestre del año 2020, en concreto, para las actividades académicas relacionadas con el tercer periodo académico del año 2020, que corresponde del xxxxxx (00) de septiembre y hasta el 19  de diciembre del 2020."}
		articulos = append(articulos, articulo)
	}
	//TCO-MTO
	if dedicacion == "TCO-MTO" {
		articulo = Articulo{Texto: "El presente acto administrativo se expide a los xxxxx (00) días del mes de XXXX del año dos mil veinte (2020) y surte efectos de conformidad con la Resolución No. 053 del veintitrés (23) de julio de 2020, por medio del cual se expide el calendario académico para el año 2020, en concreto, para las actividades académicas relacionadas en el tercer periodo académico del año 2020, comprendido entre el XXXX (XX) de septiembre al diecinueve (19) de diciembre del 2020."}
		articulos = append(articulos, articulo)
	}
	// Prestación
	if dedicacion == "HCP" && nivel == "PREGRADO" {
		articulo = Articulo{Texto: "El presente acto administrativo se expide a los XXX (XX) días del mes de XXX del año dos mil veinte (2020) y surte efectos de conformidad con la Resolución No. 053 del 23 de julio del 2020, por medio del cual se expide el calendario académico para el año 2020, en concreto, para las actividades académicas relacionadas con el tercer periodo académico del año 2020, que corresponde del xxxx (xx) de xxxxxxx y hasta 19 de diciembre del 2020.\n\n"}
		articulos = append(articulos, articulo)
	}
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "INHABILIDAD O INCOMPATIBILIDAD. - Comuníquese la presente resolución a los docentes mencionados en el artículo primero, quienes deberán manifestar, bajo la gravedad de juramento que se entiende prestado con la aprobación del formato de carga en el aplicativo de gestión académica, que no se encuentran incursos en causal de inhabilidad o incompatibilidad establecida en las normas pertinentes y aplicables, que no tienen cruces de horarios."}
		articulos = append(articulos, articulo)
	}
	/////ARTICULO 10
	//TCO-MTO
	if dedicacion == "HCP" && nivel == "POSGRADO" {
		articulo = Articulo{Texto: "VIGENCIA. -  El presente acto administrativo se expide a los XXX (XXX) días del mes de XXX del año dos mil XXX (20XX) y surte efectos de conformidad con la Resolución No. 039 de junio treinta (30) de 2020, por medio del cual se expide el calendario académico para el año 2020, en concreto, para las actividades académicas relacionadas con el tercer periodo académico del año 2020, que corresponde del xxxxxx (00) de agosto y hasta 23 de diciembre del 2020."}
		articulos = append(articulos, articulo)
	}
	resolucion.Articulos = articulos

	return resolucion
}
