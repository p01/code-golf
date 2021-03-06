package routes

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// FIXME Handle 404
func scores(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	hole := ps[0].Value
	lang := ps[1].Value

	// TODO Find a better way to do this.
	switch hole {
	case "all",
		"99-bottles-of-beer",
		"arabic-to-roman",
		"emirp-numbers",
		"evil-numbers",
		"fibonacci",
		"fizz-buzz",
		"happy-numbers",
		"odious-numbers",
		"pascals-triangle",
		"pernicious-numbers",
		"prime-numbers",
		"quine",
		"roman-to-arabic",
		"seven-segment",
		"sierpiński-triangle",
		"spelling-numbers",
		"π",
		"φ",
		"𝑒":
	default:
		print404(w, r)
		return
	}

	// TODO Find a better way to do this.
	switch lang {
	case "all", "bash", "javascript", "perl", "perl6", "php", "python", "ruby":
	default:
		print404(w, r)
		return
	}

	userID := printHeader(w, r, 200)

	w.Write([]byte(
		"<script async src=" + scoresJsPath +
			"></script><article id=scores><select id=hole>",
	))

	for _, v := range [][]string{
		{"all", "All Holes"},
		{"99-bottles-of-beer", "99 Bottles of Beer"},
		{"arabic-to-roman", "Arabic to Roman"},
		{"emirp-numbers", "Emirp Numbers"},
		{"evil-numbers", "Evil Numbers"},
		{"fibonacci", "Fibonacci"},
		{"fizz-buzz", "Fizz Buzz"},
		{"happy-numbers", "Happy Numbers"},
		{"odious-numbers", "Odious Numbers"},
		{"pascals-triangle", "Pascal's Triangle"},
		{"pernicious-numbers", "Pernicious Numbers"},
		{"prime-numbers", "Prime Numbers"},
		{"quine", "Quine"},
		{"roman-to-arabic", "Roman to Arabic"},
		{"sierpiński-triangle", "Sierpiński Triangle"},
		{"seven-segment", "Seven Segment"},
		{"spelling-numbers", "Spelling Numbers"},
		{"π", "π"},
		{"φ", "φ"},
		{"𝑒", "𝑒"},
	} {
		w.Write([]byte("<option "))
		if hole == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><select id=lang>"))

	for _, v := range [][]string{
		{"all", "All Langs"},
		{"bash", "Bash"},
		{"javascript", "JavaScript"},
		{"perl", "Perl"},
		{"perl6", "Perl 6"},
		{"php", "PHP"},
		{"python", "Python"},
		{"ruby", "Ruby"},
	} {
		w.Write([]byte("<option "))
		if lang == v[0] {
			w.Write([]byte("selected "))
		}
		w.Write([]byte("value=" + v[0] + ">" + v[1]))
	}

	w.Write([]byte("</select><table>"))

	var concat, where string

	if hole != "all" {
		where += " AND hole = '" + hole + "'"
		concat = "strokes, '<td>(', score, ' point', CASE WHEN score"
	} else {
		concat = "score, '<td>(', count, ' hole', CASE WHEN count"
	}

	if lang != "all" {
		where += " AND lang = '" + lang + "'"
	}

	rows, err := db.Query(
		`WITH leaderboard AS (
		  SELECT DISTINCT ON (hole, user_id)
		         hole,
		         submitted,
		         LENGTH(code) strokes,
		         user_id
		    FROM solutions
		   WHERE true`+where+`
		ORDER BY hole, user_id, LENGTH(code), submitted
		), scored_leaderboard AS (
		  SELECT hole,
		         ROUND(
		             (
		                 COUNT(*) OVER (PARTITION BY hole)
		                 -
		                 RANK() OVER (PARTITION BY hole ORDER BY strokes)
		                 +
		                 1
		             )
		             *
		             (
		                 100.0
		                 /
		                 COUNT(*) OVER (PARTITION BY hole)
		             )
		         ) score,
		         strokes,
		         submitted,
		         user_id
		    FROM leaderboard
		), summed_leaderboard AS (
		  SELECT user_id,
		         SUM(strokes) strokes,
		         SUM(score)   score,
		         COUNT(*),
		         MAX(submitted)
		    FROM scored_leaderboard
		GROUP BY user_id
		) SELECT CONCAT(
		             '<tr',
		             CASE WHEN user_id = $1 THEN ' class=me' END,
		             '><td>',
		             TO_CHAR(
		                 RANK() OVER (ORDER BY score DESC),
		                 'FM999"<sup>"th"</sup>"'
		             ),
		             '<td><img src="//avatars.githubusercontent.com/',
		             login,
		             '?s=26"><a href="/users/',
		             login,
		             '">',
		             login,
		             '</a><td>',
		             `+concat+` > 1 THEN 's' END,
		             ')<td>',
		             TO_CHAR(max, 'YYYY-MM-DD<span> HH24:MI:SS</span>')
		         )
		    FROM summed_leaderboard
		    JOIN users on user_id = id
		ORDER BY score DESC, max`,
		userID,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var row sql.RawBytes

		if err := rows.Scan(&row); err != nil {
			panic(err)
		}

		w.Write(row)
	}

	if err := rows.Err(); err != nil {
		panic(err)
	}
}
