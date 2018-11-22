package main

import (
	"fmt"
	"github.com/user/InterfaceGTK/lexico"
	"github.com/user/InterfaceGTK/sintactico"
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
	"os"
	"io/ioutil"
)

//-----------------------------
// Gui
//-----------------------------
type GUI struct {
	analizadorLexico *lexico.ManejadorTokens
	analizadorSintactico *sintactico.Analizador

	textViewCodigo *gtk.TextView
	textViewResultadoLexico *gtk.TextView
	textViewResultadoSintactico *gtk.TextView
	entryDireccion *gtk.Entry
	botonCargar *gtk.Button
	botonAnalizar *gtk.Button
	botonLimpiar *gtk.Button


}
func (inter *GUI) iniciar(){


	/*LEXICO*/
	anLex:=lexico.ManejadorTokens{
		Resultado:inter.textViewResultadoLexico,}//asignacion de referencias
	anLex.Construir()//carga de diccionario y vaciado de errores
	inter.analizadorLexico=&anLex

	//********************************************************
	//	Sintactico
	//********************************************************
	sintax:=sintactico.Analizador{AnalizadorLexico:inter.analizadorLexico,
		Resultado:inter.textViewResultadoSintactico,
		TextViewCodigo:inter.textViewCodigo,}//asignar referencias
	inter.analizadorSintactico=&sintax

	//--------------------------------------------------------
	// Button Cargar
	//	->leer el archivo
	//--------------------------------------------------------
	go inter.botonCargar.Clicked(func() {
		//--------------------------------------------------------
		// GtkFileChooserDialog
		//--------------------------------------------------------
		filechooserdialog := gtk.NewFileChooserDialog(
			"Choose File...",
			inter.botonCargar.GetTopLevelAsWindow(),
			gtk.FILE_CHOOSER_ACTION_OPEN,
			gtk.STOCK_OK,
			gtk.RESPONSE_ACCEPT)
		filter := gtk.NewFileFilter()
		filter.AddPattern("*.java")
		filechooserdialog.AddFilter(filter)
		filechooserdialog.Response(func() {
			//Guardar la direccion en el entry
			inter.entryDireccion.SetText(filechooserdialog.GetFilename())
			//leer el archivo
			dat,_ := ioutil.ReadFile(filechooserdialog.GetFilename())
			//check(err)

			//Insertar con buffer al text view
			buffer := inter.textViewCodigo.GetBuffer()
			//limpiar el buffer
			buffer.SetText("")
			buffer.InsertAtCursor(string(dat))

			filechooserdialog.Destroy()
		})
		filechooserdialog.Run()
	})
	//--------------------------------------------------------
	// Button Analizar
	//	->analizar sintaxis
	//--------------------------------------------------------
	inter.botonAnalizar.Clicked(func() {
		inter.analizadorSintactico.Iniciar()//vaciar errores
		inter.analizadorSintactico.Analizar()//analizar el codigo
		inter.analizadorSintactico.ImprimirResultado()
	})
	//Limpiar
	inter.botonLimpiar.Clicked(func(){
		buffer:=inter.textViewCodigo.GetBuffer()
		buffer.SetText("")
	})

	//********************************************************
	//	Cambio de Codigo
	//	->analizar Lexico
	//********************************************************
	go inter.textViewCodigo.GetBuffer().Connect("changed", func() {
		inter.analizar()
	})

}
func (inter *GUI) analizar(){
	/*	INDICES DEL BUFFER	*/
	bufferCodigo:=inter.textViewCodigo.GetBuffer()
	var start, end gtk.TextIter
	bufferCodigo.GetStartIter(&start)
	bufferCodigo.GetEndIter(&end)
	codigo:=bufferCodigo.GetText(&start, &end,false)//Obtener codigo del buffer

	inter.analizadorLexico.Iniciar(codigo)//Actualizar codigo analizado
}
func main() {
	gtk.Init(&os.Args)
	//------------------------------------------------------------
	// Window
	//------------------------------------------------------------
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetTitle("Mini JAVA")
	window.SetIconName("gtk-dialog-info")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		fmt.Println("got destroy!", ctx.Data().(string))
		gtk.MainQuit()
	}, "foo")
	//--------------------------------------------------------
	// GtkVBox
	//--------------------------------------------------------
	vbox := gtk.NewVBox(false, 1)

	//--------------------------------------------------------
	// GtkVPaned
	//--------------------------------------------------------
	vpaned := gtk.NewVPaned()
	vbox.Add(vpaned)

	//--------------------------------------------------------
	// GtkFrame
	//--------------------------------------------------------
	frame1 := gtk.NewFrame("")
	framebox1 := gtk.NewHBox(false, 1)
	frame1.Add(framebox1)

	frame2 := gtk.NewFrame("Archivo")
	framebox2 := gtk.NewVBox(false, 1)
	frame2.Add(framebox2)

	vpaned.Pack1(frame2, false, false)

	vpaned.Pack2(frame1, false, false)


	//--------------------------------------------------------
	// GtkHBox
	//--------------------------------------------------------
	buttons := gtk.NewHBox(false, 1)
	//--------------------------------------------------------
	// GtkEntry
	//--------------------------------------------------------
	entry := gtk.NewEntry()
	entry.SetEditable(false)


	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	buttonAbrir := gtk.NewButtonWithLabel("Cargar Archivo")
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	buttonAnalizar := gtk.NewButtonWithLabel("Evaluar")
	//--------------------------------------------------------
	// GtkButton
	//--------------------------------------------------------
	buttonBorrar := gtk.NewButtonWithLabel("Limpiar")


	buttons.Add(entry)
	buttons.Add(buttonAbrir)
	buttons.Add(buttonAnalizar)
	buttons.Add(buttonBorrar)
	framebox2.PackStart(buttons, false, false, 0)
	buttons = gtk.NewHBox(false, 1)

	//--------------------------------------------------------
	// GtkTextView Codigo
	//--------------------------------------------------------
	swin := gtk.NewScrolledWindow(nil, nil)
	swin.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin.SetShadowType(gtk.SHADOW_IN)
	textview := gtk.NewTextView()
	/*
	var start, end gtk.TextIter
	buffer := textview.GetBuffer()
	buffer.GetStartIter(&start)
	buffer.Insert(&start, "Error: ")
	buffer.GetEndIter(&end)
	buffer.Insert(&end, "No se ha cargado Archivo.java!")
	tag := buffer.CreateTag("bold", map[string]string{
		"background": "#FF0000", "weight": "700"})
	buffer.GetStartIter(&start)
	buffer.GetEndIter(&end)
	buffer.ApplyTag(tag, &start, &end)
	*/
	//buffer := textview.GetBuffer()
	swin.Add(textview)
	framebox2.Add(swin)


	//--------------------------------------------------------
	// GtkTextView Codigo
	//--------------------------------------------------------
	swin2 := gtk.NewScrolledWindow(nil, nil)
	swin2.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin2.SetShadowType(gtk.SHADOW_IN)
	textview2 := gtk.NewTextView()

	buffer2 := textview2.GetBuffer()
	buffer2.InsertAtCursor("\nListo para analizar")
	swin2.Add(textview2)
	frameTokns := gtk.NewFrame("\tAnálisis Léxico")
	frameTokns.Add(swin2)
	framebox1.Add(frameTokns)

	//--------------------------------------------------------
	// GtkTextView Resultado Analisis
	//--------------------------------------------------------
	swin3 := gtk.NewScrolledWindow(nil, nil)
	swin3.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	swin3.SetShadowType(gtk.SHADOW_IN)
	textview3 := gtk.NewTextView()

	buffer3 := textview3.GetBuffer()
	buffer3.InsertAtCursor("\nAnalisis")
	swin3.Add(textview3)

	frameSintactico := gtk.NewFrame("\tAnálisis Sintactixo")
	frameSintactico.Add(swin3)
	framebox1.Add(frameSintactico)


	//--------------------------------------------------------
	// Analizadores
	//--------------------------------------------------------

	//anLex.LoadDiccionario()
	//anSin:=sintactico.Analizador{}


	//--------------------------------------------------------
	// Asignación Interface GUI
	//--------------------------------------------------------
	proyeccion:=GUI{botonAnalizar:buttonAnalizar,
			textViewCodigo:textview,
			textViewResultadoLexico:textview2,
			textViewResultadoSintactico:textview3,
			entryDireccion:entry,
			botonCargar:buttonAbrir,
			botonLimpiar:buttonBorrar,
			}
	proyeccion.iniciar()
	//--------------------------------------------------------
	// Event
	//--------------------------------------------------------
	window.Add(vbox)
	window.SetSizeRequest(800, 600)
	window.ShowAll()
	gtk.Main()
}

