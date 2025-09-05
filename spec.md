# Dokumentacja Techniczna / Kontekst dla LLM: Narzędzie `getctx`

## 1. Cel i Główna Funkcjonalność

`getctx` (Get Context) to narzędzie wiersza poleceń (CLI) napisane w Go. Jego głównym celem jest **interaktywne wybieranie plików i folderów** z systemu plików, a następnie **łączenie zawartości** wszystkich wybranych (i kwalifikujących się) plików tekstowych w jeden duży plik wyjściowy (domyślnie `context.txt`).

Głównym przypadkiem użycia jest szybkie agregowanie kontekstu kodu źródłowego z projektu, który można następnie wkleić do narzędzi AI (LLM), zgłoszeń błędów lub dokumentacji.

## 2. Architektura Projektu

Projekt jest podzielony na kilka plików, z których każdy ma jasno zdefiniowaną odpowiedzialność, zgodnie z zasadą **Separacji Odpowiedzialności (Separation of Concerns)**:

- **`main.go`**: **Punkt wejścia (Entrypoint)**. Odpowiada za:

  - Parsowanie flag linii komend (np. `-o` dla pliku wyjściowego).
  - Inicjalizację i uruchomienie interfejsu TUI.
  - Przekazanie wyników z TUI do logiki budowania kontekstu.
  - Obsługę błędów na najwyższym poziomie.

- **`tui.go`**: **Terminalowy Interfejs Użytkownika (TUI)**.

  - Cała logika interaktywna jest tutaj.
  - Zbudowany w oparciu o bibliotekę **`bubbletea`** i wzorzec architektoniczny **Model-View-Update**.
  - **Model (`model` struct):** Przechowuje stan interfejsu: bieżącą ścieżkę, listę plików/folderów, pozycję kursora i mapę zaznaczonych elementów.
  - **View (`View()` method):** Renderuje stan modelu na ekranie terminala, używając stylów z biblioteki `lipgloss`.
  - **Update (`Update()` method):** Obsługuje wszystkie akcje użytkownika (naciśnięcia klawiszy) i modyfikuje stan modelu.

- **`context_builder.go`**: **Logika Biznesowa (Business Logic)**.

  - Zawiera "mózg" aplikacji, który działa po zamknięciu interfejsu TUI.
  - Odpowiada za przetworzenie listy zaznaczonych ścieżek, odfiltrowanie plików binarnych i zbudowanie finalnego pliku `context.txt`.
  - Wypisuje na konsolę szczegółowe logi dotyczące tego, które pliki są dodawane, a które pomijane.

- **`fs_utils.go`**: **Narzędzia Systemu Plików (File System Utilities)**.

  - Zbiór uniwersalnych funkcji pomocniczych do operacji na plikach i folderach.
  - **`discoverFiles`**: Rekurencyjnie przeszukuje podane ścieżki i zwraca listę wszystkich znalezionych plików, sklasyfikowanych jako tekstowe lub binarne.
  - **`isTextFile`**: Wykrywa, czy dany plik jest plikiem tekstowym (używając heurystyki opartej na typie MIME).

- **`ui.go`**: **Definicje Wyglądu (UI Definitions)**.

  - Centralne miejsce do zarządzania wyglądem aplikacji.
  - Definiuje ikony (emoji), kolory i złożone style (`lipgloss.Style`) używane w TUI i w logach.
  - Style są zorganizowane w zagnieżdżonych strukturach (`TUIStyles`, `TUIListStyles`, `TUILogStyles`) dla lepszej czytelności i skalowalności.

- **`keybindings.go`**: **Definicje Klawiszy (Keybindings)**.
  - Centralne miejsce definiujące wszystkie skróty klawiszowe używane w aplikacji jako stałe (`const`).
  - Eliminuje "magiczne stringi" w logice TUI i ułatwia rekonfigurację klawiszy.

## 3. Kluczowe Funkcjonalności i Logika

- **Nawigacja:** Użytkownik nawiguje po systemie plików za pomocą strzałek góra/dół. `Enter` wchodzi do folderu, `Backspace` cofa do folderu nadrzędnego.
- **Zaznaczanie:**
  - `Spacja`: Zaznacza/odznacza pojedynczy plik lub folder pod kursorem.
  - `CTRL+A`: Zaznacza/odznacza wszystkie elementy w bieżącym widoku.
- **Filtrowanie Plików Binarnych:** Narzędzie automatycznie wykrywa i pomija pliki binarne (np. `.exe`, obrazy, skompilowane pliki), dołączając do finalnego pliku tylko zawartość plików tekstowych. Użytkownik jest informowany o pominiętych plikach.
- **Wyjście z Programu:**
  - `q`: Kończy działanie interfejsu i **uruchamia proces budowania** pliku `context.txt` z zaznaczonych elementów.
  - `CTRL+C`: **Anuluje operację**. Kończy działanie programu bez zapisywania pliku.
- **Stylowanie:** Aplikacja intensywnie korzysta z biblioteki `lipgloss` do stylowania (kolory, pogrubienia) zarówno w interaktywnym TUI, jak i w logach wyjściowych, zapewniając spójny i nowoczesny wygląd.

## 4. Zależności Zewnętrzne

- `github.com/charmbracelet/bubbletea`: Podstawa interfejsu TUI.
- `github.com/charmbracelet/lipgloss`: Do stylowania tekstu w terminalu.
