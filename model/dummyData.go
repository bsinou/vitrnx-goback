package model

const (
	msg1 = `
Nous avons profité de ces longs mois d'hiver pour nous repencher sur notre fenêtre sur le monde (virtuel), et avons décidé de faire peau neuve.

Cette page est la première étape de cette transformation et nous vous demandons un peu de patience;

Un pas après l'autre...
`

	msg2 = `
**A lire avec délectation :**

_Fendre l'armure_ d'Anna Gavalda

Un livre plein de tendresse, d'humour et de sagesse
`
	msg3 = `
# Markdown Sheatsheet

Ci-dessous un petit récapitulatif des styles possibles à l'aide de la syntaxe "Markdown"

## a H2 Title

### a  H3 Title

## Emphasize

_italic?_ 

*Italic?* 

**bold** : 

__And here__ : 

~~strike through~~ : 

## Lists

### Unnumbered

- first
- second
- third

### Numbered

1. first
1. second
1. third 
`

	msg4 = `
Témoignage d'**Anne-Dauphine et Loïc Julliand**,
Anne-Dauphine est l'auteure d'un livre, narration de la maladie et de la mort de sa fille Thaïs atteinte d'une maladie génétique.
Anne-Dauphine est aussi la réalisatrice du film _"Et les Mistrals gagnants"_.

## Laisser la parole aux enfants.

Ils restent eux-mêmes, quelles soient les circonstances

> _"Quand on est malade, ça n'est pas difficile"_

> _"Pour vous, je sais, c'est difficile, pour moi c'est pas difficile, il faut juste patienter_

Les enfants replacent leur maladie à sa juste place

Les enfants aiment la vie en toutes circonstances

Les enfants remettent la bousssole de la vie au bon endroit en utilisant des mots simples
> _"C'est pas grave, on laisse tomber les choses qui nous tracassent et on vit avec"_

> _"Moi ça va, c'est mon coeur qui ne va pas !"_  
   Ambre, 11 ans

> _"Dans la vie il ne faut jamais se plaindre. Je suis spécial, je suis joyeux ça fait du bien parfois de pleurer.
> Je sais, on peut mourir si il n'y a pas de solution."_  
   Imad 8 ans`

	msg5 = `

- **Avoir une parole impeccable** : Cultiver attention, bienveillance sans jugement ou critique vis à vis de soi-même ou des autres.  
   Etre attentif à ces pensées, ces paroles critiques et condamnantes qui nous sont tellement habituelles qu'on n'y prête plus attention... Combien de fois par jour, par heure, par minute les commentaires que je fais dans ma tête sont-ils désagréables, insultants, critiques pour moi ou pour ceux qui m'entourent ou pour ceux à qui je pense ?

   Décider de cesser ce langage. Nettoyer ses pensées.  
   Laisser un espace à la bienveillance, à l'imagination, à l'erreur, à la découverte....

- **Ne rien prendre à titre personnel** : Quand l'autre parle, il parle de lui, même si c'est à mon propos...

- **Ne pas faire de suppositions** : Exercice bien délicat !

   Je vois ceci, j'entends cela donc....

   L'autre m'a dit.... cela signifie donc que .... et dans ma tête je crée un film, au scénario souvent dramatique.

   J'ai peur que mon conjoint, mon enfant, un de mes proches, tombe malade ou qui lui arrive un malheur... et j'écris un scénario catastrophe.... Loin de la réalité et inutile. Mais qui m'angoisse.
   Ainsi cette dame âgée de 90 ans qui toute sa vie durant a craint que l'avion que son fils prenait chaque semaine ait un accident et qui s'est angoissée des jours et des jours à ce propos, sans que cela n'arrive.... Vaines et délétères suppositions.

   Ne pas supposer. Revenir à la réalité. Oser poser des questions et ainsi obtenir des réponses.

- **Faire de son mieux** : Juste faire de son mieux

   Si je suis fatiguée après un repas, faire de mon mieux sera peut-être d'aller me reposer plutôt que de me forcer à ranger la cuisine.

   Faire de son mieux peut être dire Non à un proche qui demande un service.


Chacun de ces accords est une invitation à prendre conscience, à prêter attention au quotidien que j'encombre volontiers de lois et d'habitudes parfois obsolètes qui obscurcissent mon environnement et sont sources d'angoisse et de stress inutile.

Accepter d'en prendre conscience est un premier pas pour changer sa vie au quotidien

PAS à PAS, un PAS APRES L'AUTRE !
`
)

func populateFakeRepo() *map[string]*Article {

	tmp := make(map[string]*Article)

	tmp["new-website"] = &Article{
		Id:     "new-website",
		Tags:   "Actualités",
		Author: "Bruno",
		Title:  "Vivement le printemps !",
		Desc:   "Un court message pour vous informer que le site est en court de restructuration",
		Body:   msg1,
	}

	tmp["temps-des-vacances"] = &Article{
		Id:     "temps-des-vacances",
		Tags:   "Actualités",
		Author: "Marie-Madeleine",
		Title:  "Le temps des vacances.....",
		Desc:   "Se laisser aller à faire rien... Pourquoi pas ?",
		Body:   msg2,
	}

	tmp["body-syntax"] = &Article{
		Id:     "body-syntax",
		Tags:   "Réflexions",
		Author: "Bruno",
		Title:  "Markdown ou comment donner un style à ses textes",
		Desc:   "Un petit rappel de la syntaxe des blogs.",
		Body:   msg3,
	}

	tmp["lecon-de-vie"] = &Article{
		Id:     "lecon-de-vie",
		Tags:   "Réflexions",
		Author: "Marie-Madeleine",
		Title:  "Leçons de vie, Juillet 2017",
		Desc:   "\"Quand on ne peut pas ajouter des jours à la vie on peut toujours ajouter de la vie aux jours.\" \nPr Bernard",
		Body:   msg4,
	}

	tmp["accords-tolteques"] = &Article{
		Id:     "accords-tolteques",
		Tags:   "Réflexions",
		Author: "Marie-Madeleine",
		Title:  "Les Accords Toltèques, le 21 octobre 2016",
		Desc:   "A méditer encore et encore, les 4 accords Toltèques (Don Miguel Ruiz)",
		Body:   msg5,
	}

	return &tmp
}
