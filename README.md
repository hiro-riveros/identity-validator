# 👤 go-idvalidator

**go-idvalidator** es una librería escrita 100% en Go que permite analizar imágenes de documentos de identidad (DNI) y videos de rostros para validar la identidad de una persona. Diseñada como base para un sistema más robusto, esta librería busca ofrecer un primer nivel de coincidencia entre el documento físico y el rostro mediante procesamiento de imágenes.

---

## 🚀 Objetivo

Crear un wrapper/librería en Go que reciba como entrada:

- 2 imágenes del DNI (`frontal` y `posterior`)
- 1 video corto (o secuencia de imágenes) del rostro actual

Y entregue como salida:

```json
{
  "match_percentage": 92.5,
  "valid": true,
  "reason": "Coincidencia facial aceptable"
}
```

# 🧱 Tecnologías utilizadas
Go 1.21+

* Tesseract OCR para extraer texto del DNI
* go-face para análisis facial (bindings de Dlib)
* ffmpeg para extracción de frames desde el video
* Go image libraries para procesamiento

```bash
identity-validator/
├── go.mod
├── go.sum
├── README.md
├── main.go
├── internal/
│   ├── validator/
│   │   ├── validator.go       # Punto de entrada del análisis
│   │   ├── face_match.go      # Lógica para comparación facial
│   │   ├── dni_parser.go      # OCR y parsing del DNI
│   │   └── types.go           # Definición de structs: Input, Resultado, etc.
│   └── utils/
│       └── crypto.go          # Funciones para encriptar los datos validados
└── testdata/
    ├── sample_front_dni.jpg
    ├── sample_back_dni.jpg
    └── sample_video_face.mp4
```

# 📥 Instalación
Requisitos previos
* Instalar Tesseract:

```bash
brew install tesseract
```
* Instalar ffmpeg:
```bash
brew install ffmpeg jpeg libpng libtiff
```

* Dlib instalado (v19.10 o superior)
```bash
brew install dlib boost

# Modelos preentrenados de Dlib
wget http://dlib.net/files/shape_predictor_5_face_landmarks.dat.bz2
wget http://dlib.net/files/dlib_face_recognition_resnet_model_v1.dat.bz2
bunzip2 *.bz2
```


* Instalar librerías necesarias para go-face:

```bash
# Requiere CMake, Dlib y pkg-config
brew install cmake pkg-config
```

# ⚙️ Uso básico
Desde código Go

```go
input := validator.Input{
  FrontDNIPath: "testdata/dni_front.jpg",
  BackDNIPath:  "testdata/dni_back.jpg",
  VideoPath:    "testdata/face_video.mp4",
}

result, err := validator.Analyze(input)
if err != nil {
  log.Fatal(err)
}

fmt.Printf("Resultado: %+v\n", result)

```

Desde CLI
```bash
go run cmd/validate/main.go --front dni_front.jpg --back dni_back.jpg --video face_video.mp4
```

# 📤 Respuesta esperada
```json
{
  "match_percentage": 94.3,
  "valid": true,
  "dni_text": {
    "name": "JUAN PEREZ",
    "id_number": "12345678-9"
  },
  "reason": "Coincidencia facial aceptable"
}
```

# 🧪 Roadmap inicial
- [x] OCR de cara frontal del DNI
- [x] Extracción de frame facial del video
- [x] Matching facial con go-face
- [ ] Validación de estructura de DNI
- [ ] Persistencia opcional encriptada de resultados
- [ ] Integración futura con servicios externos

# 🔐 Seguridad
Los datos procesados se consideran sensibles. En futuras versiones se incluirá:
* Encriptación AES de los resultados
* Firma digital de los resultados
* TTL de almacenamiento temporal

# 📄 Licencia
MIT